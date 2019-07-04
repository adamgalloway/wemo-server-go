package main
 
import (
    "fmt"
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

func setupHandler(w http.ResponseWriter, r *http.Request) {
    res := setupResponse("testswitch", "aa993f4a-375f-4cf6-98d7-17bfc0f2290d", "000002F0101C00")
    w.Header().Set("Content-Type", "text/xml")
    fmt.Fprintf(w, res)
}

func upnpHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "upnp")
}


func eventHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "event")
}

func HandleHttp(port int) {
    server := http.NewServeMux()
    server.HandleFunc("/setup.xml", setupHandler)
    server.HandleFunc("/upnp/control/basicevent1", upnpHandler)
    server.HandleFunc("/eventservice.xml", eventHandler)
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), server))
}

