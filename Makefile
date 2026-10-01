.PHONY: test run broker-up broker-down firmware

test:
	go test ./...

run:
	go run ./cmd/server

broker-up:
	docker compose up -d

broker-down:
	docker compose down

firmware:
	tinygo flash -target=xiao-esp32c3 -ldflags="-X main.ssid=$${WIFI_SSID} -X main.password=$${WIFI_PASSWORD} -X main.broker=$${MQTT_BROKER}" ./firmware/esp32-c3
