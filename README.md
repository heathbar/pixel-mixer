# pixel-mixer
Open Pixel Control Mixer

pixel-mixer allows you to run multiple [OPC](http://openpixelcontrol.org/) generation scripts and fade between them. It is controlled via [MQTT](http://mqtt.org/) messages and conveniently works well as a [Home Assistant](https://home-assistant.io/) [MQTT light](https://home-assistant.io/components/light.mqtt/).

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

### Configuration Sections 
`pixel-count` is simply the number of pixels that should be expected to be received and sent from pixel-mixer.

The `mqtt` section defines how to connect to the MQTT server that will control pixel-mixer and the topics that will drive various functionality.

The `inputs` section defines which inputs should be available on the mixer. 

The `opc` section defines the size and target of the output. Typically this would be a fadecandy server.

### Extremely basic example configuration
Without defining any inputs, the mixer can only fade between solid colors using the built-in solid color generator
```JavaScript
{
    "pixel-count": 30,
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
        "destination-server": "localhost:7890"
    }
}
```
Using this simple configuration pixel-mixer has limitted functionality
```bash
# mosquitto_pub is a utility from the mosquitto-clients package.
# mosquitto_pub is not required, but useful for working with MQTT.

# RGB colors can be set via MQTT message to the `pixel-mixer/color` topic:
mosquitto_pub -h mqtt.example.com -t pixel-mixer/color -m "255,0,105"

# Output can be disabled/enabled using the `pixel-mixer/switch` topic:
mosquitto_pub -h mqtt.example.com -t pixel-mixer/switch -m "OFF"
mosquitto_pub -h mqtt.example.com -t pixel-mixer/switch -m "ON"
```


### Example configuration for two OPC inputs
This example assumes you're running pixel generation scripts from the [openpixelcontrol](https://github.com/zestyping/openpixelcontrol) python examples folder.
```JavaScript
{
    "pixel-count": 30,
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
        "destination-server": "localhost:7890"
    }
}
```
After starting pixel-mixer, start your pixel generation scripts.
```bash
raver_plaid.py localhost:7891 &
conway.py -s localhost:7892 &
```

These example scripts will now connect to pixel-mixer instead of directly to your OPC device. Using the MQTT topic and messages defined in the config, you can now fade between the two inputs as follows:

```bash
# mosquitto_pub is a utility from the mosquitto-clients package.
# mosquitto_pub is not required, but useful for working with MQTT.

# Fade to port 7891 where raver_plaid is sending pixels
mosquitto_pub -h mqtt.example.com -t pixel-mixer/input -m "raver_plaid"

# Fade to port 7892 where conway is sending pixels
mosquitto_pub -h mqtt.example.com -t pixel-mixer/input -m "conway"

# Fade to solid blue
mosquitto_pub -h mqtt.example.com -t pixel-mixer/color -m "0,0,255"

# Disable/Enable output
mosquitto_pub -h mqtt.example.com -t pixel-mixer/switch -m "OFF"
mosquitto_pub -h mqtt.example.com -t pixel-mixer/switch -m "ON"
```
### Special Inputs
If an external pixel generator is not available, pixel-mixer includes a few built-in inputs to test your setup. 

#### Solid Color
RGB colors can be set via MQTT message to the MQTT topic specified in configuration under mqtt > topics > color:
```bash
mosquitto_pub -h mqtt.example.com -t pixel-mixer/color -m "0,25,205"
```
#### Channel Walk
This input can be enable by first including it in the config:
```bash

```JavaScript
{
    "inputs":[
        {
            "type": "channel-walk",
            "mqtt-message": "channel-walk"
        }
    ]
}
```
To activate it, send the appropriate MQTT message.
```bash
mosquitto_pub -h mqtt.example.com -t pixel-mixer/input -m "channel-walk"
```
#### Rainbow
This input can be enable by first including it in the config:
```bash

```JavaScript
{
    "inputs":[
        {
            "type": "rainbow",
            "mqtt-message": "rainbow"
        }
    ]
}
```
To activate it, send the appropriate MQTT message.
```bash
mosquitto_pub -h mqtt.example.com -t pixel-mixer/input -m "rainbow"
```

