include .env
export

export PROJECT_ROOT=${shell pwd}


env-up: 
	@docker compose up -d postgres 

env-down: 
	@docker compose down postgres 


env-port-forward:
	@docker compose up -d port-forwarder


env-port-close:
	@docker compose down port-forwarder

db-wait:
	@echo "Ожидание готовности базы данных..."
	@until docker compose exec -T postgres pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}; do \
		echo "База еще не готова, ждем 1 сек..."; \
		sleep 1; \
	done;
	@echo "✅ База готова к работе!"


migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр 'seq'. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"  && \
	sudo chown $$(id -u):$$(id -g) ${PROJECT_ROOT}/migrations/*.sql && \
	echo "✅ Миграция создана и права скорректированы"

	

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр 'action'. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

migrate-up:
	@make migrate-action action=up || true


migrate-down:
	@make migrate-action action=down

	
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