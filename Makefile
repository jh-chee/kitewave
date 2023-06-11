pre:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/cloudwego/kitex/tool/cmd/kitex@latest
	go install github.com/cloudwego/thriftgo@latest

generate:
	mkdir -p ./http-server/proto_gen
	protoc -I=. --go_out=./http-server/proto_gen ./idl_http.proto
	cd http-server && kitex -module github.com/jh-chee/kitewave/http-server ../idl_rpc.thrift
	cd rpc-server && kitex -module github.com/jh-chee/kitewave/rpc-server ../idl_rpc.thrift

test:
	go test -cover ./http-server ./rpc-server -coverprofile=coverage.out

lint:
	golangci-lint run -v --fast \
		--disable-all \
		--enable=gofmt \
		--enable=govet \
		--enable=staticcheck \
		--enable=unused \
		--enable=structcheck \
		--enable=errcheck \
		./http-server/... ./rpc-server/...

build:
	go build -v ./...

mock:
	cd http-server && mockery --keeptree -r --all
	cd rpc-server && mockery --keeptree -r --all

compose:
	docker-compose build kitewave-http-server 
	docker-compose build kitewave-rpc-server 
	docker-compose up

docker-build: docker-build-http docker-build-rpc

docker-build-http:
	docker build ./http-server -t kitewave-http-server:latest

docker-build-rpc:
	docker build ./rpc-server -t kitewave-rpc-server:latest

.PHONY: pre generate test lint build mock compose docker-build docker-build-http docker-build-rpc
