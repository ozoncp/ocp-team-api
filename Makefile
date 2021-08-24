LOCAL_BIN:=$(CURDIR)/bin

run:
	go run cmd/ocp-team-api/main.go

lint:
	golint ./...

test:
	go test -v ./...

.PHONY: build
build: vendor-proto .generate .build

PHONY: .generate
.generate:
		mkdir -p swagger
		mkdir -p pkg/ocp-team-api
		protoc -I vendor.protogen \
				--go_out=pkg/ocp-team-api --go_opt=paths=import \
				--go-grpc_out=pkg/ocp-team-api --go-grpc_opt=paths=import \
				--grpc-gateway_out=pkg/ocp-team-api \
				--grpc-gateway_opt=logtostderr=true \
				--grpc-gateway_opt=paths=import \
				--swagger_out=allow_merge=true,merge_file_name=api:swagger \
				api/ocp-team-api/ocp-team-api.proto
		mv pkg/ocp-team-api/github.com/ozoncp/ocp-team-api/pkg/ocp-team-api/* pkg/ocp-team-api/
		rm -rf pkg/ocp-team-api/github.com
		mkdir -p cmd/ocp-team-api
		cd pkg/ocp-team-api && ls go.mod || go mod init github.com/ozoncp/ocp-team-api/pkg/ocp-team-api && go mod tidy

.PHONY: generate
generate: .vendor-proto .generate

.PHONY: .build
.build:
		go build -o bin/ocp-team-api cmd/ocp-team-api/main.go

.PHONY: vendor-proto
vendor-proto: .vendor-proto

.PHONY: .vendor-proto
.vendor-proto:
		mkdir -p vendor.protogen
		mkdir -p vendor.protogen/api/ocp-team-api
		cp api/ocp-team-api/ocp-team-api.proto vendor.protogen/api/ocp-team-api/ocp-team-api.proto
		@if [ ! -d vendor.protogen/google ]; then \
			git clone https://github.com/googleapis/googleapis vendor.protogen/googleapis &&\
			mkdir -p  vendor.protogen/google/ &&\
			mv vendor.protogen/googleapis/google/api vendor.protogen/google &&\
			rm -rf vendor.protogen/googleapis ;\
		fi


.PHONY: deps
deps: install-go-deps

.PHONY: install-go-deps
install-go-deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
		ls go.mod || go mod init github.com/ozoncp/ocp-team-api
		go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
		go get -u github.com/golang/protobuf/proto
		go get -u github.com/golang/protobuf/protoc-gen-go
		go get -u google.golang.org/grpc
		go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go get -u github.com/envoyproxy/protoc-gen-validate
		go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
		go install github.com/envoyproxy/protoc-gen-validate


.PHONY: migrate
migrate: .install-migrate-deps .migrate

.PHONY: .install-migrate-deps
.install-migrate-deps:
		go get -u github.com/pressly/goose/v3/cmd/goose

.PHONY: .migrate
.migrate:
		goose -dir migrations postgres "postgres://root:root@localhost:5432/postgres?sslmode=disable" up