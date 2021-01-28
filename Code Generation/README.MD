To compile proto files with grpc plugin you will need to install ***protoc*** compiler and ***grpc plugin*** for golang. 

1. Finish the prerequisites section from the following link:

	***https://grpc.io/docs/languages/go/quickstart/***

2. After we will need to download grpc plugin with the following command:
 
	***go get -u google.golang.org/grpc***

This will download grpc pluning at $GOPATH/src/google.golang.org/grpc. 

3. To build those protofiles located on source folder to greetpb output folder run make build command.

	***make build***