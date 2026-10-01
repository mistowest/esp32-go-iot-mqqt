# ESP32-C3 firmware

This firmware is written in Go and compiled with TinyGo. It connects an ESP32-C3 to Wi-Fi, reads a DHT22 and publishes an MQTT JSON payload.

## Wiring

| DHT22 | XIAO ESP32-C3 |
|---|---|
| VCC | 3.3V |
| DATA | GPIO4 / D2 |
| GND | GND |

Use a 4.7k–10k pull-up from DATA to 3.3V if your DHT22 module does not already contain one.

## Build

Install TinyGo 0.42+ and use a board target with Wi-Fi support such as `xiao-esp32c3` or `esp32c3-generic`. Native Wi-Fi on these ESP32 variants is provided by `tinygo.org/x/espradio`.

Set credentials and broker with linker variables:

```bash
tinygo flash -target=xiao-esp32c3 \
  -ldflags='-X main.ssid=YOUR_WIFI -X main.password=YOUR_PASSWORD -X main.broker=192.168.1.100:1883' \
  ./firmware/esp32-c3
```

Change `dhtPin` in `main.go` if your DHT22 DATA wire uses another GPIO.

The DHT22 should not be sampled faster than once every two seconds; this firmware uses a three-second loop.
