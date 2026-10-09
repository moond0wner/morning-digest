include .env
export

run:
	@go run cmd/main.go

go-lint:
	@~/go/bin/golangci-lint run --config=.golangci.yml ./...

logs-cleanup:
	@read -p "Очистить все log файлы окружения? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		sudo rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi;