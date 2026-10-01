//go:build tinygo

package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	mqtt "github.com/soypat/natiu-mqtt"
	"machine"
	"tinygo.org/x/drivers/dht"
	"tinygo.org/x/drivers/netdev"
	nl "tinygo.org/x/drivers/netlink"
	link "tinygo.org/x/espradio/netlink"
)

var (
	ssid     string
	password string
	broker   = "192.168.1.100:1883"
	deviceID = "esp32-dht22-01"
	topic    = "iot/dht22/esp32-dht22-01/data"
)

const dhtPin = machine.GPIO4

func main() {
	time.Sleep(2 * time.Second)
	if err := connectWiFi(); err != nil { panic(err) }
	sensor := dht.NewDummyDevice(dhtPin, dht.DHT22)
	for {
		if err := publishReading(sensor); err != nil {
			fmt.Println("MQTT error:", err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		time.Sleep(3 * time.Second)
	}
}

func connectWiFi() error {
	radio := link.Esplink{}
	netdev.UseNetdev(&radio)
	for attempt := 1; attempt <= 5; attempt++ {
		fmt.Printf("WiFi: connecting (attempt %d)\n", attempt)
		err := radio.NetConnect(&nl.ConnectParams{Ssid: ssid, Passphrase: password})
		if err == nil { fmt.Println("WiFi: connected"); return nil }
		fmt.Println("WiFi:", err.Error())
		time.Sleep(3 * time.Second)
	}
	return errors.New("could not connect to WiFi")
}

func publishReading(sensor dht.DummyDevice) error {
	if err := sensor.ReadMeasurements(); err != nil { return fmt.Errorf("DHT22 read: %w", err) }
	temperature, err := sensor.TemperatureFloat(dht.C)
	if err != nil { return err }
	humidity, err := sensor.HumidityFloat()
	if err != nil { return err }
	payload := fmt.Sprintf(`{"device_id":"%s","temperature":%.1f,"humidity":%.1f,"timestamp":%d}`, deviceID, temperature, humidity, time.Now().Unix())
	conn, err := net.Dial("tcp", broker)
	if err != nil { return fmt.Errorf("broker connection: %w", err) }
	defer conn.Close()
	client := mqtt.NewClient(mqtt.ClientConfig{Decoder: mqtt.DecoderNoAlloc{make([]byte, 1500)}})
	var connect mqtt.VariablesConnect
	connect.SetDefaultMQTT([]byte(deviceID))
	connect.KeepAlive = 30
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx, conn, &connect)
	cancel()
	if err != nil { return fmt.Errorf("MQTT connect: %w", err) }
	flags, err := mqtt.NewPublishFlags(mqtt.QoS0, false, false)
	if err != nil { return err }
	variables := mqtt.VariablesPublish{TopicName: []byte(topic)}
	if err := client.PublishPayload(flags, variables, []byte(payload)); err != nil { return fmt.Errorf("MQTT publish: %w", err) }
	fmt.Printf("published temperature=%.1f°C humidity=%.1f%%\n", temperature, humidity)
	client.Disconnect(nil)
	return nil
}
