package main
 
import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
    "strconv"
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

func upnpHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var state string
        body, err := ioutil.ReadAll(r.Body)
        if err == nil && strings.Contains(string(body), "<BinaryState>1</BinaryState>") {
            state = "1"
            // turn on
        } else {
            state = "0"
            // turn off
        }
        res := upnpResponse(state)
        w.Header().Set("Content-Type", "text/xml")
        fmt.Fprintf(w, res)
    }
}


func eventHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "event")
}

func HandleHttp(port int) {
    server := http.NewServeMux()
    server.HandleFunc("/setup.xml", setupHandler("testswitch", "aa993f4a-375f-4cf6-98d7-17bfc0f2290d", "000002F0101C00"))
    server.HandleFunc("/upnp/control/basicevent1", upnpHandler())
    server.HandleFunc("/eventservice.xml", eventHandler)
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), server))
}

