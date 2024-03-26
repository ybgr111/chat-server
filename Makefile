LOCAL_BIN:=$(CURDIR)/bin

##linter
install-golangci-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

lint:
	$(LOCAL_BIN)/golangci-lint run ./... --config .golangci.pipeline.yaml

##proto; grpc
install-deps:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.14.0

get-deps:
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc


generate:
	make generate-note-api

generate-note-api:
	mkdir -p pkg/note_v1
	protoc --proto_path api/note_v1 \
	--go_out=pkg/note_v1 --go_opt=paths=source_relative \
	--plugin=protoc-gen-go=bin/protoc-gen-go \
	--go-grpc_out=pkg/note_v1 --go-grpc_opt=paths=source_relative \
	--plugin=protoc-gen-go-grpc=bin/protoc-gen-go-grpc \
	api/note_v1/note.proto

##build
build:
	GOOS=linux GOARCH=amd64 go build -o service_linux cmd/grpc_server/main.go

copy-to-server:
	scp service_linux root@185.91.52.223:

docker-build-and-push:
	docker buildx build --no-cache --platform linux/amd64 -t cr.selcloud.ru/ybgr111/chat-server:v0.0.1 .
	docker login -u token -p CRgAAAAAmIxM5SY6qVc7pMYCMUQDKkRmX0KBpU3A cr.selcloud.ru/ybgr111
	docker push cr.selcloud.ru/ybgr111/chat-server:v0.0.1

##migrations
local-migration-status:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir ${LOCAL_MIGRATION_DIR} postgres ${LOCAL_MIGRATION_DSN} down -v

##tests
.PHONY: test
test:
	go clean -testcache
	go test ./... -covermode count -coverpkg=github.com/ybgr111/chat-server/internal/service/...,github.com/ybgr111/chat-server/internal/api/... -count 5

.PHONY: test-coverage
test-coverage:
	go clean -testcache
	go test ./... -coverprofile=coverage.tmp.out -covermode count -coverpkg=github.com/ybgr111/chat-server/internal/service/...,github.com/ybgr111/chat-server/internal/api/... -count 5
	grep -v 'mocks\|config' coverage.tmp.out  > coverage.out
	rm coverage.tmp.out
	go tool cover -html=coverage.out;
	go tool cover -func=./coverage.out | grep "total";
	grep -sqFx "/coverage.out" .gitignore || echo "/coverage.out" >> .gitignore