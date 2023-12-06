//build protobuf

.PHONY: api
api:
	protoc -I api/ \
		-I${GOPATH}/src \
		-I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
		--go_out=plugins=grpc:api \
		api/api.proto