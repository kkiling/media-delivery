TEST_DB_NAME:=media_delivery_test
TEST_DB_USERNAME:=root
TEST_DB_PASSWORD:=root
DB_HOST:=localhost
DB_PORT:=5432
GOPATH_BIN := $(GOPATH)/bin

.PHONY: bin-deps
bin-deps:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/pav5000/smartimports/cmd/smartimports@v0.2.0
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/matm/gocov-html/cmd/gocov-html@latest
	go install github.com/axw/gocov/gocov@latest
	go install github.com/ofabry/go-callvis@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install tool

.PHONY: mocks
mocks:
	@echo "\n --- 🤡 Create Mocks --- \n"
	go generate ./...

.PHONY: test
test:
	@echo "\n --- 🧪 Run project tests --- \n"
	go test ./...

coverage:
	go test -coverprofile=coverage.out ./...
	gocov convert coverage.out | gocov-html > coverage.html
	open coverage.html 2>/dev/null || xdg-open coverage.html 2>/dev/null || echo "Open coverage.html manually"
	@echo "🖲️ Coverage summary:"
	@go tool cover -func=coverage.out | grep total | awk '{printf "\033[32m✅ Total coverage: %s\033[0m\n", $$3}'

.PHONY: format
format:
	@echo "\n --- 🚀 Start format imports --- \n"
	smartimports -local "github.com/kkiling/media-delivery/" -path . #-exclude pkg/mocks

.PHONY: test-db
test-db:
	@echo "\n --- 🖲️ Migrate test postgresql database --- \n"
	@echo "\n --- 🖲️ Dropping and creating test database --- \n"
	# Удаляем и создаем базу данных используя правильные учетные данные
	psql postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres -c "DROP DATABASE IF EXISTS ${TEST_DB_NAME};"
	psql postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/postgres -c "CREATE DATABASE ${TEST_DB_NAME};"
	@echo "\n --- 🖲️ Applying  postgresql migrations --- \n"
	# Накатываем миграции на тестовую базу
	goose -dir=./migrations/postgresql/state postgres "postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${TEST_DB_NAME}?sslmode=disable" up
	goose -dir=./migrations/postgresql postgres "postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${TEST_DB_NAME}?sslmode=disable" up

	@echo "\n --- 🖲️ Creating .testenv file --- \n"
	rm -f ./.testenv
	# Создаем файл с строкой подключения для PostgreSQL
	echo "POSTGRES_CONN_STRING=postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${TEST_DB_NAME}?sslmode=disable" > ./.testenv

	@echo "\n --- ✅ Test database setup completed --- \n"

.PHONY: schema
schema:
	pg_dump "postgresql://${TEST_DB_USERNAME}:${TEST_DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${TEST_DB_NAME}?sslmode=disable" --no-owner --no-privileges --no-tablespaces --no-security-labels --no-comments -s >  schema.sql
	sed -i '/^\\restrict/d' schema.sql
	sed -i '/^\\unrestrict/d' schema.sql
	sqlc generate

clear-cache:
	buf registry cc

.PHONY: generate
generate:
	mkdir -p vendor.protogen
	cp -R api/media-delivery/ vendor.protogen/
	buf dep update
	buf generate

.PHONY: run
run:
	go run ./cmd

build-docker:
	docker build -f build/Dockerfile -t media-delivery-server:latest .

# Граф проекта
.PHONY: callvis
callvis:
	go-callvis -focus=main -http=:7878 ./cmd