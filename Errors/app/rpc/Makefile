build: 
	rm -rf languagepb
	mkdir languagepb
	protoc -I. proto/*.proto --go_out=languagepb --go-grpc_out=languagepb

.PHONY: build