package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mistowest/esp32-go-iot-mqqt/internal/config"
)

type Reading struct {
	DeviceID    string  `json:"device_id"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Timestamp   int64   `json:"timestamp"`
}

func main() {
	cfg := config.Load()

	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(5 * time.Second)
	opts.SetKeepAlive(30)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
		opts.SetPassword(cfg.Password)
	}

	opts.OnConnect = func(client mqtt.Client) {
		log.Printf("[MQTT] connected to %s", cfg.Broker)
		if token := client.Subscribe(cfg.Topic, cfg.QoS, onMessage); token.Wait() && token.Error() != nil {
			log.Printf("[MQTT] subscribe error: %v", token.Error())
			return
		}
		log.Printf("[MQTT] subscribed to %s", cfg.Topic)
	}

	opts.OnConnectionLost = func(_ mqtt.Client, err error) {
		log.Printf("[MQTT] connection lost: %v", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("[MQTT] connection failed: %v", token.Error())
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down")
	client.Disconnect(250)
}

func onMessage(_ mqtt.Client, message mqtt.Message) {
	var reading Reading
	if err := json.Unmarshal(message.Payload(), &reading); err != nil {
		log.Printf("[DATA] invalid JSON on %s: %v", message.Topic(), err)
		return
	}

	stamp := time.Unix(reading.Timestamp, 0).Format(time.RFC3339)
	log.Printf("[DATA] device=%s temperature=%.1f°C humidity=%.1f%% timestamp=%s",
		reading.DeviceID, reading.Temperature, reading.Humidity, stamp)
}
