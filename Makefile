.PHONY: dev_setup dev_fresh setup recover fresh update license version

dev_setup:
	@./scripts/install-xray-mac.sh

dev_run:
	go run main.go start

dev_fresh:
	rm -f storage/app/*.json
	rm -f storage/app/*.txt
	rm -f storage/database/*.json
	rm -f storage/logs/*.log

setup:
	./scripts/setup-updater.sh
	@if [ ! -f ./configs/main.local.json ]; then \
		cp ./configs/main.json ./configs/main.local.json; \
	fi

recover:
	@./scripts/recover.sh

fresh:
	rm -f storage/app/*.json
	rm -f storage/app/*.txt
	rm -f storage/database/*.json
	rm -f storage/logs/*.log
	docker compose restart

update: setup
	@git pull
	@./scripts/update.sh

license:
	@./scripts/license.sh "$(v)"

version:
	@docker compose exec app ./xray-manager version
