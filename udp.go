package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"text/template"
)

func lookupHost() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	CheckError(err)
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	i := strings.LastIndex(localAddr, ":")
	return localAddr[0:i]
}

var searchResponse = template.New("searchResponse")

func init() {

	deviceResponse := []string{
		"HTTP/1.1 200 OK",
		"ST: urn:Belkin:device:**",
		"USN: uuid:{{.id}}::urn:Belkin:device:controllee:1",
		"LOCATION: http://{{.host}}:{{.port}}/setup.xml",
		"CACHE-CONTROL: max-age=1800",
		"SERVER: Unspecified, UPnP/1.0, Unspecified",
		"OPT: \"http://schemas.upnp.org/upnp/1/0/\"; ns=01",
		"01-NLS: {{.id}}",
		"X-User-Agent: redsonic",
	}

	searchResponse.Parse(strings.Join(deviceResponse, "\r\n") + "\r\n\r\n")
}

//HandleUDP create udp handlers for each device
func HandleUDP(devices map[string]Device) {
	serverAddr, err := net.ResolveUDPAddr("udp", "239.255.255.250:1900")
	CheckError(err)

	serverConn, err := net.ListenMulticastUDP("udp", nil, serverAddr)
	CheckError(err)
	defer serverConn.Close()

	buf := make([]byte, 1024)

	host := lookupHost()
	fmt.Println("Udp ip", host)

	for {
		n, addr, err := serverConn.ReadFromUDP(buf)
		if err != nil {
			log.Println("Error reading from UDP:", err)
		}

		req := string(buf[0:n])

		if strings.Contains(req, "M-SEARCH") &&
			(strings.Contains(req, "urn:Belkin:device:**") ||
				strings.Contains(req, "ssdp:all") ||
				strings.Contains(req, "upnp:rootdevice")) {

			fmt.Println("Received belkin upnp from:", addr)
			fmt.Println("Request:", req)

			// loop over devices
			for _, device := range devices {
				conn, err := net.Dial("udp", addr.String())
				if err == nil {
					port := strconv.Itoa(device.Port)
					id := device.ID
					var tpl bytes.Buffer
					searchResponse.Execute(&tpl, map[string]string{"host": host, "port": port, "id": id})
					conn.Write(tpl.Bytes())
					conn.Close()

					fmt.Println("Responded to:", addr)
				}
			}
		}
	}
}
