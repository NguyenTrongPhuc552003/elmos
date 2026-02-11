# Makefile for generating gRPC code from proto files

PROTO_DIR := api/proto
PROTO_FILES := $(wildcard $(PROTO_DIR)/*.proto)

.PHONY: all proto clean

all: proto

proto:
	@echo "Generating Go code from proto files..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
	@echo "✓ Generated gRPC code successfully"

clean:
	@rm -f $(PROTO_DIR)/*.pb.go
	@echo "✓ Cleaned generated files"
