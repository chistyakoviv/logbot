DOCKER_COMPOSE_FILE ?= compose.yaml
IMAGE_TAG ?= latest

# Project commands
migrate-status:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot-go-cli make logbot-migrate-status

migrate-up:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot-go-cli make logbot-migrate-up

migrate-down:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot-go-cli make logbot-migrate-down

migration:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot-go-cli goose create ${MIGRATION_NAME} sql

lint:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot-go-cli make logbot-lint

lint-migrate:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot-go-cli golangci-lint migrate --config .golangci.pipeline.yaml

.PHONY: mocks
mocks:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot-go-cli mockery

.PHONY: tests
tests:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot go test -v ./...

db-connect:
	docker compose -f ${DOCKER_COMPOSE_FILE} exec logbot-pg psql postgres://app:secret@logbot-pg/app

db-purge:
	docker compose -f ${DOCKER_COMPOSE_FILE} run --rm logbot-go-cli make logbot-db-purge

wait-db:
	wait-for-it logbot-pg:5432 -t 60

# DB commands
logbot-migrate-status: wait-db
	goose postgres "${POSTGRES_DSN}" status -v

logbot-migrate-up: wait-db
	goose postgres "${POSTGRES_DSN}" up -v

logbot-migrate-down: wait-db
	goose postgres "${POSTGRES_DSN}" down -v

logbot-db-purge: wait-db
	psql ${POSTGRES_DSN} -t -c "SELECT 'DROP TABLE \"' || tablename || '\" CASCADE;' FROM pg_tables WHERE schemaname = 'public'" | \
	psql ${POSTGRES_DSN}

logbot-lint:
# Scan directories separately to avoid golangci-lint access errors caused by directories owned by root
	golangci-lint run -v ./internal/... ./cmd/... --config .golangci.pipeline.yaml

# Git commands
version:
# tag v1.0.1
	git tag "${TAG}"
	git push origin "${TAG}"

domain:
	ngrok http 80

# Docker commands
build: build-logbot

build-logbot:
	docker --log-level=info buildx build \
	--pull \
	--target builder \
	--cache-from type=registry,ref=${REGISTRY}/logbot-builder:buildcache\
	--cache-to type=registry,ref=${REGISTRY}/logbot-builder:buildcache,mode=max \
	--tag ${REGISTRY}/logbot-builder:${IMAGE_TAG}\
	--file docker/production/logbot/Dockerfile \
	--push \
	./

	docker --log-level=info buildx build \
	--cache-from type=registry,ref=${REGISTRY}/logbot:buildcache \
	--cache-to type=registry,ref=${REGISTRY}/logbot:buildcache,mode=max \
	--tag ${REGISTRY}/logbot:${IMAGE_TAG} \
	--file docker/production/logbot/Dockerfile \
	--push \
	./

	docker --log-level=info buildx build \
	--cache-from type=registry,ref=${REGISTRY}/logbot-go-cli:buildcache \
	--cache-to type=registry,ref=${REGISTRY}/logbot-go-cli:buildcache,mode=max \
	--tag ${REGISTRY}/logbot-go-cli:${IMAGE_TAG} \
	--file docker/production/go-cli/Dockerfile \
	--push \
	./

pull: pull-logbot

pull-logbot:
	docker pull ${REGISTRY}/logbot:${IMAGE_TAG}
	docker pull ${REGISTRY}/logbot-go-cli:${IMAGE_TAG}

tag: tag-logbot

tag-logbot:
	docker tag ${REGISTRY}/logbot:${SRC_IMAGE_TAG} ${REGISTRY}/logbot:${IMAGE_TAG}
	docker tag ${REGISTRY}/logbot-go-cli:${SRC_IMAGE_TAG} ${REGISTRY}/logbot-go-cli:${IMAGE_TAG}

push: push-logbot

push-logbot:
	docker push ${REGISTRY}/logbot:${IMAGE_TAG}
	docker push ${REGISTRY}/logbot-go-cli:${IMAGE_TAG}
	docker push ${REGISTRY}/logbot:${IMAGE_TAG}-builder

push-latest: push-logbot-latest

push-logbot-latest:
	docker push ${REGISTRY}/logbot:${IMAGE_TAG}
	docker push ${REGISTRY}/logbot-go-cli:${IMAGE_TAG}
