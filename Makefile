MIGRATE=migrate -path sql/migrations -database postgres://postgres:mysecretpassword@localhost:5432/postgres?sslmode=disable

# ==============================================================================
# Building containers
all: all lint test taskmanager

taskmanager:
	docker build \
		-f deployments/dockerfile.tasktracker-api \
		-t tasktracker-api-amd64:1.0 \
		--build-arg VCS_REF=`git rev-parse HEAD` \
		--build-arg BUILD_DATE=`date -u +”%Y-%m-%dT%H:%M:%SZ”` \
		.

# ==============================================================================
# Lint all files in the project

lint:
	@golangci-lint run -c .golangci.yml

# ==============================================================================
# Running tests within the local computer

test:
	go test ./... -count=1

# ==============================================================================
# Running from within docker compose

run: up 

up:
	docker-compose -f deployments/compose/compose.yaml  up --detach --remove-orphans

down:
	docker-compose -f deployments/compose/compose.yaml down --remove-orphans

logs:
	docker-compose -f deployments/compose/compose.yaml logs -f

# ==============================================================================
# Modules support

deps-reset:
	git checkout -- go.mod
	go mod tidy
	go mod vendor

tidy:
	go mod tidy
	go mod vendor

deps-upgrade:
	# go get $(go list -f '{{if not (or .Main .Indirect)}}{{.Path}}{{end}}' -m all)
	go get -u -t -d -v ./...
	go mod tidy
	go mod vendor

deps-cleancache:
	go clean -modcache

# ==============================================================================
# Docker support

FILES := $(shell docker ps -aq)

down-local:
	docker stop $(FILES)
	docker rm $(FILES)

clean:
	docker system prune -f	

logs-local:
	docker logs -f $(FILES)

# ==============================================================================
# Migration support

migrate-create: 
	migrate create -ext sql -dir sql/migrations migration

migrate-up: ## Run migrations
	$(MIGRATE) up

migrate-down: ## Rollback migrations
	$(MIGRATE) down
