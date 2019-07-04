package main
 
import (
    "fmt"
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

func searchResponse(host string, port int, id string) string {
    var res strings.Builder
    res.WriteString("HTTP/1.1 200 OK\r\n")
    res.WriteString("CACHE-CONTROL: max-age=86400\r\n")
    res.WriteString("DATE: Sat, 26 Nov 2016 04:56:29 GMT\r\n")
    res.WriteString("EXT:\r\n")
    res.WriteString("LOCATION: http://")
    res.WriteString(host)
    res.WriteString(":")
    res.WriteString(strconv.Itoa(port))
    res.WriteString("/setup.xml\r\n")
    res.WriteString("OPT: \"http://schemas.upnp.org/upnp/1/0/\"; ns=01\r\n")
    res.WriteString("01-NLS: b9200ebb-736d-4b93-bf03-835149d13983\r\n")
    res.WriteString("SERVER: Unspecified, UPnP/1.0, Unspecified\r\n")
    res.WriteString("ST: urn:Belkin:device:**\r\n")
    res.WriteString("USN: uuid:")
    res.WriteString(id)
    res.WriteString("::urn:Belkin:device:**\r\n")
    res.WriteString("X-User-Agent: redsonic\r\n\r\n")
    return res.String()
}
 
func HandleUpnp(devices map[string]Device) {
    serverAddr,err := net.ResolveUDPAddr("udp","239.255.255.250:1900")
    CheckError(err)
 
    serverConn,err := net.ListenMulticastUDP("udp", nil, serverAddr)
    CheckError(err)
    defer serverConn.Close()
 
    buf := make([]byte, 1024)

    host := lookupHost()
    fmt.Println(host)
 
    for {
        n,addr,err := serverConn.ReadFromUDP(buf)
        if err != nil {
            fmt.Println("Error: ",err)
        } 

        req := string(buf[0:n])
 
        if strings.Contains(req, "M-SEARCH") && 
            (strings.Contains(req, "urn:Belkin:device:**") || 
            strings.Contains(req, "ssdp:all") || 
            strings.Contains(req, "upnp:rootdevice")) {
 
           fmt.Println("Received belkin upnp request")

            // loop over devices
            for _, device := range devices {
                conn,err := net.Dial("udp", addr.String())
                if err == nil {
                    fmt.Fprintf(conn, searchResponse(host, device.Port, device.Id))
                    conn.Close()
                }
            }
        }
    }
}

