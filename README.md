# pixel-mixer
Open Pixel Control Mixer

pixel-mixer allows you to run multiple [OPC](http://openpixelcontrol.org/) generation scripts and fade between them. It is controlled via [MQTT](http://mqtt.org/) messages and not-coincidentally works well as a [Home Assistant](https://home-assistant.io/) [MQTT light](https://home-assistant.io/components/light.mqtt/).

## Setup
```bash
git clone https://github.com/heathbar/pixel-mixer.git
cd pixel-mixer
go get github.com/kellydunn/go-opc
go get github.com/eclipse/paho.mqtt.golang
go build
./pixel-mixer
```

## Configuration
A configuration file can be specified with the `-c config.file.json` argument or `./config.json` is used by default.

The `mqtt` section defines how to connect to the MQTT server that will control pixel-mixer and the topics that will drive various functionality.

The `inputs` section defines which inputs should be available on the mixer. 

The `opc` section defines the size and target of the output. Typically this would be a fadecandy server.

### Extremely basic example configuration
Without defining any inputs, the mixer can only fade bwtween solid colors using the built-in solid color generator
```JavaScript
{
    "mqtt": {
        "server": "tcp://mqtt.example.com:1883",
        "topics": {
            "power": "pixel-mixer/switch",
            "input": "",
            "color": "pixel-mixer/color"
        }
    },
    "inputs":[],
    "opc": {
        "destination-server": "localhost:7890",
        "pixel-count": 30
    }
}
```
With this configuration you can send RGB colors via MQTT message to the `pixel-mixer/color` topic. For example:
```bash
# mosquitto_pub is a utility from the mosquitto-clients package.
# mosquitto_pub is not required, but useful for working with MQTT.
mosquitto_pub -h mqtt.example.com -t pixel-mixer/color -m "255,0,105"
```

### Example configuration for two OPC inputs
This example assumes you're running pixel generation scripts from the [openpixelcontrol](https://github.com/zestyping/openpixelcontrol) examples.
```JavaScript
{
    "mqtt": {
        "server": "tcp://mqtt.example.com:1883",
        "topics": {
            "power": "pixel-mixer/switch",
            "input": "pixel-mixer/input",
            "color": "pixel-mixer/color"
        }
    },
    "inputs":[
        {
            "type": "opc",
            "mqtt-message": "raver_plaid"
            "port": 7891
        },
        {
            "type": "opc",
            "mqtt-message": "conway"
            "port": 7892
        }
    ],
    "opc": {
        "destination-server": "localhost:7890",
        "pixel-count": 30
    }
}
```
After starting pixel-mixer, start your pixel generation scripts.
```bash
raver_plaid.py localhost:7891 &
conway.py -s localhost:7892 &
```

These example scripts will now connect to pixel-mixer instead of directly to your OPC device. Using the mqtt topic and messages defined in the config, you can now fade between the two inputs as follows:

```bash
# mosquitto_pub is a utility from the mosquitto-clients package.
# mosquitto_pub is not required, but useful for working with MQTT.
mosquitto_pub -h mqtt.example.com -t pixel-mixer/input -m "raver_plaid"
# or
mosquitto_pub -h mqtt.example.com -t pixel-mixer/input -m "conway"
```