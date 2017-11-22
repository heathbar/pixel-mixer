# pixel-mixer
Open Pixel Control Mixer

pixel-mixer allows you to run multiple [OPC](http://openpixelcontrol.org/) generation scripts and fade between them. It is controlled via [MQTT](http://mqtt.org/) messages and not-coincidentally works well as a [Home Assistant](https://home-assistant.io/) [MQTT light](https://home-assistant.io/components/light.mqtt/).

### Setup
```bash
git clone https://github.com/heathbar/pixel-mixer.git
cd pixel-mixer
go get github.com/kellydunn/go-opc
go get github.com/eclipse/paho.mqtt.golang
go build
./pixel-mixer
```
