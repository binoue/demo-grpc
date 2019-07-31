PROTOBUF_VERSION=3.9.0
CURL=curl -Lsf

all: generate

bin/protoc:
	mkdir -p bin
	rm -rf include
	$(CURL) -o protobuf.zip https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOBUF_VERSION)/protoc-$(PROTOBUF_VERSION)-linux-x86_64.zip
	unzip protobuf.zip bin/protoc 'include/*'
	rm -f protobuf.zip

bin/protoc-gen-go:
	GOBIN="$(shell pwd)/bin" go get -u github.com/golang/protobuf/protoc-gen-go

proto/timereport.pb.go: timereport.proto bin/protoc bin/protoc-gen-go
	mkdir -p proto
	PATH="$(shell pwd)/bin:$(PATH)" bin/protoc -I. --go_out=plugins=grpc,paths=source_relative:proto $<

generate: proto/timereport.pb.go

.PHONY: all generate build
