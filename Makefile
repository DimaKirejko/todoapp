include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "CLEANUP ALL VOLUME?? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		rm -rf out/pgdata && \
		echo "DB file is cleared"; \
	else \
		echo "Attempt to clear DB rejected"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "seq is required. example: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		create -ext sql -dir /migrations -seq "$(seq)"

migrate-up:
	@ make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "action is required. example: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database "postgres://${PG_USER}:${PG_PASS}@todoapp-postgres:5432/${PG_DB}?sslmode=disable" \
		$(action)

display:
	@echo "$(PG_USER) $(PG_PASS) $(PG_DB)"

run-todoapp:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export PG_HOST=localhost && \
	go mod tidy && \
	go run cmd/todoapp/main.go

# inst: 7:48