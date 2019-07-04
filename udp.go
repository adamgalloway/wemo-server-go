package main
 
import (
    "fmt"
    "net"
    "net/http"
    "os"
    "strings"
    "strconv"
    "log"
)

func SearchResponse(host string, port int, id string) string {
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
 
/* A Simple function to verify error */
func CheckError(err error) {
    if err  != nil {
        fmt.Println("Error: " , err)
        os.Exit(0)
    }
}
 
func HandleUpnp() {
    /* Lets prepare a address at any address at port 1900*/   
    ServerAddr,err := net.ResolveUDPAddr("udp","239.255.255.250:1900")
    CheckError(err)
 
    /* Now listen at selected port */
    ServerConn, err := net.ListenMulticastUDP("udp", nil, ServerAddr)
    CheckError(err)
    defer ServerConn.Close()
 
    buf := make([]byte, 1024)
 
    for {
        n,addr,err := ServerConn.ReadFromUDP(buf)
        req := string(buf[0:n])

        fmt.Println("Received ",req," from ",addr)
 
        if strings.Contains(req, "M-SEARCH") && 
           (strings.Contains(req, "urn:Belkin:device:**") || 
           strings.Contains(req, "ssdp:all") || 
           strings.Contains(req, "upnp:rootdevice")) {
            fmt.Println("Received belkin request")
        }

        if err != nil {
            fmt.Println("Error: ",err)
        } 
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func HandleHttp(port int) {
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}

func main() {
    go HandleHttp(8080)
    HandleUpnp()
}
