package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
)

var upnpResponse = template.New("upnpResponse")
var setupResponse = template.New("setupResponse")

func init() {
	upnpResponse.Parse(`<?xml version="1.0"?>` +
		`<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">` +
		`<s:Body>` +
		`<u:{{.method}}BinaryStateResponse xmlns:u="urn:Belkin:service:basicevent:1">` +
		`<BinaryState>{{.state}}</BinaryState>` +
		`</u:{{.method}}BinaryStateResponse>` +
		`</s:Body>` +
		`</s:Envelope>`)

	setupResponse.Parse(`<?xml version="1.0"?>` +
		`<root xmlns="urn:Belkin:device-1-0">` +
		`<specVersion><major>1</major><minor>0</minor></specVersion>` +
		`<device>` +
		`<deviceType>urn:Belkin:device:controllee:1</deviceType>` +
		`<binaryState>1</binaryState>` +
		`<friendlyName>{{.name}}</friendlyName>` +
		`<manufacturer>Belkin International Inc.</manufacturer>` +
		`<manufacturerURL>http://www.belkin.com</manufacturerURL>` +
		`<modelDescription>Belkin Plugin Socket 1.0</modelDescription>` +
		`<modelName>Socket</modelName>` +
		`<modelNumber>1</modelNumber>` +
		`<modelURL>http://www.belkin.com/plugin/</modelURL>` +
		`<serialNumber>{{.serial}}</serialNumber>` +
		`<UDN>uuid:{{.id}}</UDN>` +
		`<UPC>123456789</UPC>` +
		`<serviceList>` +
		`<service>` +
		`<serviceType>urn:Belkin:service:basicevent:1</serviceType>` +
		`<serviceId>urn:Belkin:serviceId:basicevent1</serviceId>` +
		`<controlURL>/upnp/control/basicevent1</controlURL>` +
		`<eventSubURL>/upnp/event/basicevent1</eventSubURL>` +
		`<SCPDURL>/eventservice.xml</SCPDURL>` +
		`</service>` +
		`</serviceList>` +
		`</device>` +
		`</root>`)
}

func setupHandler(name string, id string, serial string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Setup request from", r.RemoteAddr)

		w.Header().Set("CONTENT-TYPE", "text/xml")
		setupResponse.Execute(w, map[string]string{"name": name, "id": id, "serial": serial})
	}
}

func upnpHandler(oncommand string, offcommand string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Upnp request from", r.RemoteAddr)

		var command string
		var method string = "Get"
		var state string = "0"

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body")
		}

		bodyString := string(body)

		if strings.Contains(bodyString, "SetBinaryState") {
			// TODO return state
			method = "Set"

			if strings.Contains(bodyString, "<BinaryState>1</BinaryState>") {
				fmt.Println("Received upnp on")
				state = "1"
				// turn on
				command = oncommand
			} else if strings.Contains(bodyString, "<BinaryState>0</BinaryState>") {
				fmt.Println("Received upnp off")
				state = "0"
				// turn off
				command = offcommand
			}
		}

		w.Header().Set("CONTENT-TYPE", "text/xml")
		upnpResponse.Execute(w, map[string]string{"method": method, "state": state})

		if len(command) > 0 {
			go runCommand(command)
		}
	}
}

func runCommand(command string) {
	fmt.Println("Executing command:", command)
	cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("Execution error:", err)
	} else {
		fmt.Println("Execution result:", out.String())
	}
}

var eventResponse string = `<?xml version="1.0"?>
<scpd xmlns="urn:Belkin:service-1-0">
<actionList>
  <action>
    <name>SetBinaryState</name>
    <argumentList>
      <argument>
        <retval/>
        <name>BinaryState</name>
        <relatedStateVariable>BinaryState</relatedStateVariable>
        <direction>in</direction>
      </argument>
    </argumentList>
  </action>
  <action>
    <name>GetBinaryState</name>
    <argumentList>
      <argument>
        <retval/>
        <name>BinaryState</name>
        <relatedStateVariable>BinaryState</relatedStateVariable>
        <direction>out</direction>
      </argument>
    </argumentList>
  </action>
</actionList>
<serviceStateTable>
  <stateVariable sendEvents="yes">
    <name>BinaryState</name>
    <dataType>Boolean</dataType>
    <defaultValue>0</defaultValue>
  </stateVariable>
  <stateVariable sendEvents="yes">
    <name>level</name>
    <dataType>string</dataType>
    <defaultValue>0</defaultValue>
  </stateVariable>
</serviceStateTable>
</scpd>`

func eventHandler(w http.ResponseWriter, r *http.Request) {
	res := eventResponse
	w.Header().Set("CONTENT-TYPE", "text/xml")
	fmt.Fprintf(w, res)
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	w.WriteHeader(404)
	fmt.Fprintf(w, "not found")
}

func handleHTTP(name string, device Device) {
	fmt.Println("Starting server on", device.Port)
	server := http.NewServeMux()
	server.HandleFunc("/", logHandler)
	server.HandleFunc("/setup.xml", setupHandler(name, device.Id, device.Serial))
	server.HandleFunc("/upnp/control/basicevent1", upnpHandler(device.OnCommand, device.OffCommand))
	server.HandleFunc("/eventservice.xml", eventHandler)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(device.Port), server))
}

//HandleHTTP create http handlers for each device
func HandleHTTP(devices map[string]Device) {

	for key, device := range devices {
		go handleHTTP(key, device)
	}

}
