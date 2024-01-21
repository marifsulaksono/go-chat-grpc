GO_WORKSPACE := ..

protoc:
	protoc --experimental_allow_proto3_optional --go_out=$(GO_WORKSPACE) --go-grpc_out=require_unimplemented_servers=false:$(GO_WORKSPACE) --proto_path=protos protos/*.proto 
	@echo "protoc compile done!"