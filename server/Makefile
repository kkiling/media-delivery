TEST_DB_NAME:=./media-delivery.db
TEST_DB_PATH := $(CURDIR)/$(notdir ${TEST_DB_NAME})
GOPATH_BIN := $(GOPATH)/bin

.PHONY: bin-deps
bin-deps:
	go install github.com/golang/mock/mockgen@v1.6.0
	go install github.com/pav5000/smartimports/cmd/smartimports@v0.2.0
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install tool

.PHONY: mocks
mocks:
	@echo "\n --- 🤡 Create Mocks --- \n"
	go generate ./...

.PHONY: test
test:
	@echo "\n --- 🧪 Run project tests --- \n"
	go test ./...

.PHONY: format
format:
	@echo "\n --- 🚀 Start format imports --- \n"
	smartimports -local "github.com/kkiling/media-delivery/" -path . #-exclude pkg/mocks

.PHONY: test-db
test-db:
	@echo "\n --- 🖲️ Migrate test sqlite database --- \n"
	rm -f ${TEST_DB_PATH}
	rm -f "${TEST_DB_PATH}-shm"
	rm -f "${TEST_DB_PATH}-wal"
	@echo "\n --- 🖲️ statemachine sqlite migrations --- \n"
	goose  -dir=./migrations/sqlite/state sqlite3 ${TEST_DB_PATH} up
	@echo "\n --- 🖲️ media-delivery sqlite migrations --- \n"
	goose -dir=./migrations/sqlite sqlite3 ${TEST_DB_PATH} up
	@echo "\n --- 🖲️ Creating .testenv file --- \n"
	rm -f ./testenv
	echo "SQLITE_DSN=${TEST_DB_PATH}" > ./.testenv

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
