package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//Device struct to hold device definition
type Device struct {
	ID         string `json:"uuid"`
	Serial     string `json:"serial"`
	OnCommand  string `json:"oncommand"`
	OffCommand string `json:"offcommand"`
	Port       int    `json:"port"`
}

//CheckError helper to log fatal errors and exit
func CheckError(err error) {
	if err != nil {
		log.Fatalln("Error: ", err)
	}
}

func loadDevices(file string) map[string]Device {
	jsonFile, err := os.Open(file)
	CheckError(err)

	fmt.Println("Opened", file)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]Device
	json.Unmarshal([]byte(byteValue), &result)

	return result
}

func main() {
	deviceFile := os.Args[1]
	devices := loadDevices(deviceFile)

	HandleHTTP(devices)
	HandleUDP(devices)
}
