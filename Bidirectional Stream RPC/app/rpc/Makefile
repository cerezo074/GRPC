build: 
	rm -rf imagepb
	mkdir imagepb
	protoc -I. proto/*.proto --go_out=imagepb --go-grpc_out=imagepb

.PHONY: build