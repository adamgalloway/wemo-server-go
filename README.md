# Emulated Wemo Plug

This project provides config file based virtual plugs for use with Amazon Alexa.

The devices.json specifies on/off commands to be controlled with Alexa.

## Setup

Setup init script (linux)

```
sudo cp wemo.sh /etc/init.d/wemo
sudo update-rc.d wemo defaults
```

## Executable Usage

Copy or build the executable (for your architecture) into the working directory used in the init script. Add your devices.json to the same directory or modify the init script with the path to the desired devices.json file.

```
cp bin/wemo-linux-arm /home/pi/wemo/wemo
cp devices.json /home/pi/wemo/
```

Build script based on: https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04
