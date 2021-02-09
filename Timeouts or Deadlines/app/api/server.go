package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"service/core"
	"service/rpc/languagepb"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LanguageHandler struct {
	languagepb.UnsafeLanguageServiceServer
}

func (*LanguageHandler) Detect(ctx context.Context, request *languagepb.LanguageDetectorRequest) (*languagepb.LanguageDetectorResponse, error) {
	fmt.Printf("Language detector function was invoked with %v", request)
	text := request.GetText()

	if strings.TrimSpace(text) == "" {
		return nil, status.Error(codes.InvalidArgument, "Text can not be empty.")
	}

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

func (handler *LanguageHandler) DetectWithDeadline(ctx context.Context, request *languagepb.LanguageDetectorRequest) (*languagepb.LanguageDetectorResponse, error) {
	someLongProcess() //Trigger a dealine or timeout, take too much time

	if ctx.Err() == context.Canceled {
		return nil, status.Error(codes.Canceled, "Client cancelled, abandoning.")
	}

	return handler.Detect(ctx, request)
}

func someLongProcess() {
	time.Sleep(time.Millisecond * 2100)
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
