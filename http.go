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
        w.Header().Set("Content-Type", "text/xml")
        fmt.Fprintf(w, res)
    }
}

func upnpResponse(state string) string {
    var res strings.Builder
    res.WriteString("<s:Envelope xmlns:s=\"http://schemas.xmlsoap.org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap.org/soap/encoding/\">")
    res.WriteString("<s:Body>")
    res.WriteString("<u:SetBinaryStateResponse xmlns:u=\"urn:Belkin:service:basicevent:1\">")
    res.WriteString("<BinaryState>")
    res.WriteString(state)
    res.WriteString("</BinaryState>")
    res.WriteString("</u:SetBinaryStateResponse>")
    res.WriteString("</s:Body>")
    res.WriteString("</s:Envelope>\r\n\r\n")
    return res.String()
}

func upnpHandler(oncommand string, offcommand string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var state string
        var command string
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
           log.Fatal(err)
        }

        bodyString := string(body)

        if strings.Contains(bodyString, "GetBinaryState") {
            // TODO return state
            return
        } else if strings.Contains(bodyString, "<BinaryState>1</BinaryState>") {
            state = "1"
            // turn on
            command = oncommand
        } else if strings.Contains(bodyString, "<BinaryState>0</BinaryState>") {
            state = "0"
            // turn off
            command = offcommand
        } else {
            return
        }

        cmd := exec.Command("sh", "-c", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
            log.Fatal(err)
	}
	fmt.Println("execution result: ", out.String())

        res := upnpResponse(state)
        w.Header().Set("Content-Type", "text/xml")
        fmt.Fprintf(w, res)
    }
}

func eventResponse() string {
    var res strings.Builder
    res.WriteString("<scpd xmlns=\"urn:Belkin:service-1-0\">")
    res.WriteString("<actionList>")
    res.WriteString("<action>")
    res.WriteString("<name>SetBinaryState</name>")
    res.WriteString("<argumentList>")
    res.WriteString("<argument>")
    res.WriteString("<retval/>")
    res.WriteString("<name>BinaryState</name>")
    res.WriteString("<relatedStateVariable>BinaryState</relatedStateVariable>")
    res.WriteString("<direction>in</direction>")
    res.WriteString("</argument>")
    res.WriteString("</argumentList>")
    res.WriteString("</action>")
    res.WriteString("</actionList>")
    res.WriteString("<serviceStateTable>")
    res.WriteString("<stateVariable sendEvents=\"yes\">")
    res.WriteString("<name>BinaryState</name>")
    res.WriteString("<dataType>Boolean</dataType>")
    res.WriteString("<defaultValue>0</defaultValue>")
    res.WriteString("</stateVariable>")
    res.WriteString("<stateVariable sendEvents=\"yes\">")
    res.WriteString("<name>level</name>")
    res.WriteString("<dataType>string</dataType>")
    res.WriteString("<defaultValue>0</defaultValue>")
    res.WriteString("</stateVariable>")
    res.WriteString("</serviceStateTable>")
    res.WriteString("</scpd>\r\n\r\n")
    return res.String()
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
    res := eventResponse()
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
