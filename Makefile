.PHONY: setup run build reset fresh

setup:
	./third_party/install-xray-mac.sh

run:
	go run main.go start

build:
	go build main.go -o xray-manager

fresh:
	rm storage/database.json
	rm storage/xray.json
	docker compose restart

update:
	docker compose pull
	git pull
	docker compose down
	docker compose up -d
