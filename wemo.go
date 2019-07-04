package main

import (
    "io/ioutil"
    "fmt"
    "os"
    "encoding/json"
)

type Device struct {
    Id string `json:"uuid"`
    Serial string `json:"serial"`
    OnCommand string `json:"oncommand"`
    OffCommand string `json:"offcommand"`
}

func LoadDevices(file string) map[string]Device {
    jsonFile, err := os.Open(file)
    CheckError(err)

    fmt.Println("Opened ",file)
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var result map[string]Device
    json.Unmarshal([]byte(byteValue), &result)

    return result
}

func main() {

    deviceFile := os.Args[1]

    devices := LoadDevices(deviceFile)
    //fmt.Println(devices)

    i := 8080
    for key, device := range devices { 
        go HandleHttp(i, key, device.Id, device.Serial, device.OnCommand, device.OffCommand)
        i++
    }

    HandleUpnp("192.168.1.34", devices)
}
