package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"service/cmd/rpc/languagepb"
	"service/core"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

type LanguageHandler struct {
	languagepb.UnsafeLanguageServiceServer
}

func (*LanguageHandler) Detect(ctx context.Context, request *languagepb.LanguageDetectorRequest) (*languagepb.LanguageDetectorResponse, error) {
	fmt.Printf("Language detector function was invoked with %v", request)
	text := request.GetText()
	result := core.NewLanguageResult(text)
	response := &languagepb.LanguageDetectorResponse{
		Result: &languagepb.Language{
			Name:       result.Name,
			Script:     result.Script,
			Confidence: result.Confidence,
		},
	}

	return response, nil
}

func main() {
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}

	fmt.Println("Loading SSL Credentials...")
	certFile := "../cmd/tls/output/server.crt"
	keyFile := "../cmd/tls/output/server.pem"
	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("Failed loading certificates: %v", sslErr)
		return
	}

	opts := grpc.Creds(creds)
	grpcServer := grpc.NewServer(opts)
	reflection.Register(grpcServer)
	languagepb.RegisterLanguageServiceServer(grpcServer, &LanguageHandler{})
	fmt.Println("Server Listening...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
