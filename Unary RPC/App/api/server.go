package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"service/core"
	"service/rpc/languagepb"

	"google.golang.org/grpc"
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

	grpcServer := grpc.NewServer()
	languagepb.RegisterLanguageServiceServer(grpcServer, &LanguageHandler{})
	fmt.Println("Server Listening...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
