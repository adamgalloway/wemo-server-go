package main
 
import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "strconv"
    "os/exec"
    "bytes"
    "log"
    "text/template"
)

func setupResponse(name string, id string, serial string) string {
    var res strings.Builder
    res.WriteString("<?xml version=\"1.0\"?>")
    res.WriteString("<root xmlns=\"urn:Belkin:device-1-0\">")
        res.WriteString("<specVersion>")
            res.WriteString("<major>1</major>")
            res.WriteString("<minor>0</minor>")
        res.WriteString("</specVersion>")
        res.WriteString("<device>")
             res.WriteString("<deviceType>urn:Belkin:device:controllee:1</deviceType>")
             res.WriteString("<friendlyName>")
             res.WriteString(name)
             res.WriteString("</friendlyName>")
             res.WriteString("<manufacturer>Belkin International Inc.</manufacturer>")
             res.WriteString("<modelName>Emulated Socket</modelName>")
             res.WriteString("<modelNumber>3.1415</modelNumber>")
             res.WriteString("<manufacturerURL>http://www.belkin.com</manufacturerURL>")
             res.WriteString("<modelDescription>Belkin Plugin Socket 1.0</modelDescription>")
             res.WriteString("<modelURL>http://www.belkin.com/plugin/</modelURL>")
             res.WriteString("<UDN>uuid:")
             res.WriteString(id)
             res.WriteString("</UDN>")
             res.WriteString("<serialNumber>")
             res.WriteString(serial)
             res.WriteString("</serialNumber>")
             res.WriteString("<binaryState>0</binaryState>")
             res.WriteString("<serviceList>")
                 res.WriteString("<service>")
                     res.WriteString("<serviceType>urn:Belkin:service:basicevent:1</serviceType>")
                     res.WriteString("<serviceId>urn:Belkin:serviceId:basicevent1</serviceId>")
                     res.WriteString("<controlURL>/upnp/control/basicevent1</controlURL>")
                     res.WriteString("<eventSubURL>/upnp/event/basicevent1</eventSubURL>")
                     res.WriteString("<SCPDURL>/eventservice.xml</SCPDURL>")
                 res.WriteString("</service>")
             res.WriteString("</serviceList>")
         res.WriteString("</device>")
    res.WriteString("</root>\r\n\r\n")
    return res.String()
}

func setupHandler(name string, id string, serial string) http.HandlerFunc {
    res := setupResponse(name, id, serial)
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("setup request from", r.RemoteAddr)

        w.Header().Set("Content-Type", "text/xml")
        fmt.Fprintf(w, res)
    }
}

var upnpResponse = template.New("upnpResponse")

func init() {
  upnpResponse.Parse(`<?xml version="1.0"?>
  <s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
    <s:Body>
      <u:{{method}}BinaryStateResponse xmlns:u="urn:Belkin:service:basicevent:1">
        <BinaryState>{{state}}</BinaryState>
      </u:{{method}}BinaryStateResponse>
    </s:Body>
  </s:Envelope>\r\n`)
}

func upnpHandler(oncommand string, offcommand string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("upnp request from", r.RemoteAddr)

        var command string
        var method string = "Get"
        var state string ="0"

        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
           fmt.Println("error reading body")
           log.Fatal(err)
        }

        bodyString := string(body)

        if strings.Contains(bodyString, "SetBinaryState") {
            // TODO return state
            method = "Set"
        }

        if strings.Contains(bodyString, "<BinaryState>1</BinaryState>") {
            fmt.Println("on")
            state = "1"
            // turn on
            command = oncommand
        } else if strings.Contains(bodyString, "<BinaryState>0</BinaryState>") {
            fmt.Println("off")
            state = "0"
            // turn off
            command = offcommand
        }

        if len(command) > 0 {
            fmt.Println("executing command", command)
            cmd := exec.Command("sh", "-c", command)
	    var out bytes.Buffer
	    cmd.Stdout = &out
	    err = cmd.Run()
	    if err != nil {
                log.Fatal(err)
	    }
	    fmt.Println("execution result: ", out.String())
        }
        w.Header().Set("Content-Type", "text/xml")
        upnpResponse.Execute(w, map[string]string{"method": method,"state": state})
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
</scpd>\r\n`

func eventHandler(w http.ResponseWriter, r *http.Request) {
    res := eventResponse
    w.Header().Set("Content-Type", "text/xml")
    fmt.Fprintf(w, res)
}

func handleHttp(name string, device Device) {
    fmt.Println("Starting server on ", device.Port)
    server := http.NewServeMux()
    server.HandleFunc("/setup.xml", setupHandler(name, device.Id, device.Serial))
    server.HandleFunc("/upnp/control/basicevent1", upnpHandler(device.OnCommand, device.OffCommand))
    server.HandleFunc("/eventservice.xml", eventHandler)
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(device.Port), server))
}

func HandleHttp(devices map[string]Device) {

    for key, device := range devices {
        go handleHttp(key, device)
    }

}
