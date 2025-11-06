# Project commands
cli-test:
	docker compose run --rm go-cli sh -c "CONFIG_PATH=config/local.yml go run cmd/logbot/main.go"

migrate-status:
	docker compose run --rm go-cli make logbot-migrate-status

migrate-up:
	docker compose run --rm go-cli make logbot-migrate-up

migrate-down:
	docker compose run --rm go-cli make logbot-migrate-down

migration:
	docker compose run --rm go-cli goose create ${MIGRATION_NAME} sql

lint:
	docker compose run --rm go-cli make logbot-lint

lint-migrate:
	docker compose run --rm go-cli golangci-lint migrate --config .golangci.pipeline.yaml

.PHONY: mocks
mocks:
	docker compose run --rm go-cli mockery

.PHONY: tests
tests:
	docker compose run --rm logbot go test -v ./...

db-connect:
	docker compose exec logbot-pg psql postgres://app:secret@logbot-pg/app

db-purge:
	docker compose exec logbot-pg sh -c "psql postgres://app:secret@logbot-pg/app -t -c \"SELECT 'DROP TABLE \\\"' || tablename || '\\\" CASCADE;' FROM pg_tables WHERE schemaname = 'public'\" | psql postgres://app:secret@logbot-pg/app"

wait-db:
	wait-for-it logbot-pg:5432 -t 60

# DB commands
logbot-migrate-status: wait-db
	goose postgres "${PG_DSN}" status -v

logbot-migrate-up: wait-db
	goose postgres "${PG_DSN}" up -v

logbot-migrate-down: wait-db
	goose postgres "${PG_DSN}" down -v

logbot-lint:
	golangci-lint run -v ./... --config .golangci.pipeline.yaml

# Git commands
tag:
# tag v1.0.1
	git tag "${TAG}"
	git push origin "${TAG}"
