include .env
export

export PROJECT_ROOT=${shell pwd}

run:
	@go run cmd/main.go

deploy:
	@docker compose up -d --build tgbot

undeploy:
	@docker compose down tgbot

ps:
	@docker compose ps

go-lint:
	@~/go/bin/golangci-lint run --config=.golangci.yml --issues-exit-code=0 ./...

logs-cleanup:
	@read -p "Очистить все log файлы окружения? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		sudo rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi;