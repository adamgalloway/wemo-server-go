# Emulated Wemo Plug

[![Go Report Card](https://goreportcard.com/badge/github.com/adamgalloway/wemo-server-go)](https://goreportcard.com/report/github.com/adamgalloway/wemo-server-go)

This project provides config file based virtual plugs for use with Amazon Alexa.

The devices.json specifies on/off commands to be controlled with Alexa.

## Setup

Setup init script (linux)

```
sudo cp wemo.service /etc/systemd/system
sudo systemctl enable wemo.service
```

## Executable Usage

Copy or build the executable (for your architecture) from [Releases](https://github.com/adamgalloway/wemo-server-go/releases) into the working directory used in the init script. Add your devices.json to the same directory or modify the init script with the path to the desired devices.json file.

```
cp bin/wemo-linux-arm /home/pi/wemo/wemo
cp devices.json /home/pi/wemo/
```

## Device Config

The json to define the list of devices should be configured with on and off commands for each device as follows:

```
{

	"living-room-tv": {
 		"uuid" : "94de9b11-5115-46ab-88f8-dc1a7d440918",
		"serial" : "000002F0101C00",
		"oncommand"  : "/home/pi/turn-on-tv.sh",
		"offcommand" : "/home/pi/turn-off-tv.sh",
		"port" : 8080
	},

	"living-room-switch": { 
		"uuid" : "81abe8d4-a481-47d2-9def-c3c581eb0ed7",
		"serial" : "000001F0101C00",
		"oncommand"  : "/home/pi/turn-on-switch.sh",
		"offcommand" : "/home/pi/turn-off-switch.sh",
		"port" : 8081
	}

}
```

Build script based on: https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
