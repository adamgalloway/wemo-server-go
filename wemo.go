package main

import (
    "io/ioutil"
    "log"
    "os"
    "encoding/json"
)

type Device struct {
    Id string `json:"uuid"`
    Serial string `json:"serial"`
    OnCommand string `json:"oncommand"`
    OffCommand string `json:"offcommand"`
    Port int
}

func CheckError(err error) {
    if err  != nil {
        log.Fatalln("Error: " , err)
    }
}

func LoadDevices(file string) map[string]Device {
    jsonFile,err := os.Open(file)
    CheckError(err)

    log.Println("Opened",file)
    defer jsonFile.Close()

    byteValue, _ := ioutil.ReadAll(jsonFile)

    var result map[string]Device
    json.Unmarshal([]byte(byteValue), &result)

    i := 9000
    for key, device := range result {
        device.Port = i
        result[key] = device
        i++
    }

    return result
}

func main() {
    deviceFile := os.Args[1]
    devices := LoadDevices(deviceFile)

    HandleHttp(devices)
    HandleUdp(devices)
}
