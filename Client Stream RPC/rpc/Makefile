build: 
	rm -rf videopb
	mkdir videopb
	protoc -I. proto/*.proto --go_out=videopb --go-grpc_out=videopb

.PHONY: build