include .env
export 

export PROJECT_ROOT=${shell pwd}
MAKEFLAGS += --no-print-directory

env-up:
	@docker compose up -d postgres

env-down:
	@docker compose down postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres port-forwarder && \
		rm -rf ${PROJECT_ROOT}/out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "${seq}" ]; then \
		echo "Отсутствует необходимый параметр 'seq'. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir //migrations \
		-seq "${seq}"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "${action}" ]; then \
		echo "Отсутствует необходимый параметр 'action'. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		-path //migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable \
		${action}

logs-cleanup:
	@read -p "Очистить все log файлы? Опасность утери логов. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi

run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/messenger/main.go

deploy:
	@docker compose up -d --build messenger

undeploy:
	@docker compose down messenger

ps:
	@docker compose ps

env-config:
	@docker compose config


swagger-gen:
	@docker compose run --rm swagger \
		init \
		-g cmd/messenger/main.go \
		-o docs \
		--parseInternal \
		--parseDependency

diagram-db:
	@npx --yes @mermaid-js/mermaid-cli@11.16.0 \
		-i docs/database.mmd \
		-o docs/database.svg \
		-b transparent

test-unit:
	@go test ${or ${action},./...}

test-env-up:
	@- docker compose exec postgres \
		psql -U ${POSTGRES_USER} -d postgres \
		-c "CREATE DATABASE ${POSTGRES_TEST_DB};"


test-migrate-up:
	@docker compose run --rm postgres-migrate \
		-path //migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_TEST_DB}?sslmode=disable \
		up

test-env-down:
	@- docker compose exec postgres \
		psql -U $(POSTGRES_USER) -d postgres \
		-c "DROP DATABASE $(POSTGRES_TEST_DB) WITH (FORCE);"

test-integration:
	@ export POSTGRES_HOST=localhost && \
	export POSTGRES_DB=${POSTGRES_TEST_DB} && \
	go test -tags=integration -count=1 ${or ${action},./...}
