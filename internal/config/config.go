package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
	Topic    string
	QoS      byte
}

func Load() Config {
	return Config{
		Broker:   env("MQTT_BROKER", "tcp://localhost:1883"),
		ClientID: env("MQTT_CLIENT_ID", "dht22-server"),
		Username: os.Getenv("MQTT_USERNAME"),
		Password: os.Getenv("MQTT_PASSWORD"),
		Topic:    env("MQTT_TOPIC", "iot/dht22/+/data"),
		QoS:      qos(env("MQTT_QOS", "0")),
	}
}

func env(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func qos(value string) byte {
	n, err := strconv.Atoi(value)
	if err != nil || n < 0 || n > 2 {
		return 0
	}
	return byte(n)
}
