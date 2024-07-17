postgres:
	docker run --name keeper -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=Xer_0101 -d postgres:14-alpine
	sleep 10
	$(MAKE) db_create
	sleep 1
	$(MAKE) db_uuid
db_create:
	docker exec -it keeper psql -U postgres -c "CREATE DATABASE keeper;"
db_uuid:
	docker exec -it keeper psql -U postgres -d keeper -c "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"

# Linter constants
LINTER := golangci-lint
lint:
	@echo === Lint
	$(LINTER) --version
	$(LINTER) cache clean && $(LINTER) run

generate:
	${GO} generate ./...

reload_postgres:
	docker stop keeper > /dev/null
	docker rm keeper > /dev/null
	$(MAKE) postgres

test:
	go test ./...

.PHONY: db_create db_uuid postgres lint generate reload_postgres test
