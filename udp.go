package main
 
import (
    "text/template"
    "log"
    "net"
    "strings"
    "strconv"
)

func lookupHost() string {
    conn,err := net.Dial("udp", "8.8.8.8:80")
    CheckError(err)
    defer conn.Close()
    localAddr := conn.LocalAddr().String()
    i := strings.LastIndex(localAddr, ":")
    return localAddr[0:i]
}

var searchResponse = template.New("searchResponse")

func init() {
  searchResponse.Parse(`HTTP/1.1 200 OK
CACHE-CONTROL: max-age=86400
DATE: Sat, 26 Nov 2016 04:56:29 GMT
EXT:
LOCATION: http://{{.host}}:{{.port}}/setup.xml
OPT: "http://schemas.upnp.org/upnp/1/0/"; ns=01
01-NLS: b9200ebb-736d-4b93-bf03-835149d13983
SERVER: Unspecified, UPnP/1.0, Unspecified
ST: urn:Belkin:device:**
USN: uuid:{{.id}}::urn:Belkin:device:**
X-User-Agent: redsonic`)
}
 
func HandleUdp(devices map[string]Device) {
    serverAddr,err := net.ResolveUDPAddr("udp","239.255.255.250:1900")
    CheckError(err)
 
    serverConn,err := net.ListenMulticastUDP("udp", nil, serverAddr)
    CheckError(err)
    defer serverConn.Close()
 
    buf := make([]byte, 1024)

    host := lookupHost()
    log.Println("Udp ip:", host)
 
    for {
        n,addr,err := serverConn.ReadFromUDP(buf)
        if err != nil {
            log.Println("Error reading from UDP:", err)
        } 

        req := string(buf[0:n])
 
        if strings.Contains(req, "M-SEARCH") && 
            (strings.Contains(req, "urn:Belkin:device:**") || 
            strings.Contains(req, "ssdp:all") || 
            strings.Contains(req, "upnp:rootdevice")) {
 
            log.Println("Received belkin upnp from:", addr)

            // loop over devices
            for _, device := range devices {
                conn,err := net.Dial("udp", addr.String())
                if err == nil {
                    port :=  strconv.Itoa(device.Port)
                    id := device.Id
                    searchResponse.Execute(conn, map[string]string{"host": host, "port": port, "id": id})
                    conn.Close()
                }
            }
        }
    }
}

