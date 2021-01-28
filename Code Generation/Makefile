build: 
	rm -rf greetpb
	mkdir greetpb
	protoc -I. source/*.proto --go_out=greetpb --go-grpc_out=greetpb

.PHONY: build