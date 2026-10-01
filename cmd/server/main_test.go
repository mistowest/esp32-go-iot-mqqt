package main

import (
	"encoding/json"
	"testing"
)

func TestReadingJSON(t *testing.T) {
	input := `{"device_id":"esp32-dht22-01","temperature":26.4,"humidity":58.7,"timestamp":1790890200}`
	var reading Reading
	if err := json.Unmarshal([]byte(input), &reading); err != nil {
		t.Fatal(err)
	}
	if reading.DeviceID != "esp32-dht22-01" {
		t.Fatalf("unexpected device: %q", reading.DeviceID)
	}
	if reading.Temperature != 26.4 || reading.Humidity != 58.7 {
		t.Fatalf("unexpected values: %+v", reading)
	}
}
