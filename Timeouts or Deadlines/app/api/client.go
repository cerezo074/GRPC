package main

import (
	"context"
	"service/rpc/languagepb"
	"time"

	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Running Client...")

	clientConnection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
		return
	}

	defer clientConnection.Close()

	connection := languagepb.NewLanguageServiceClient(clientConnection)
	run(connection, "上海＝南部さやか】世界保健機関（ＷＨＯ）の国際調査団が本格調査を始めた中国湖北省武漢市で、新型コロナウイルスの感染で家族を亡くした遺族に対し、当局からの圧力が強まっている。")
	runWithDeadline(connection, "Esta variante se propaga más rápido y podría conllevar mayor riesgo de infección, por eso, los investigadores creen que esta nueva variante se propagó de Brasil a Asia.")
}

func run(connection languagepb.LanguageServiceClient, text string) {
	request := &languagepb.LanguageDetectorRequest{
		Text: text,
	}
	response, err := connection.Detect(context.Background(), request)
	if err != nil {
		statusError, ok := status.FromError(err)
		if ok {
			log.Printf("Error detected, type: %v, description: %v", statusError.Code(), statusError.Message())
			return
		}

		log.Println("Error while calling Detect RPC: %v", err)
		return
	}

	log.Printf("Response from Detect Language Service.\n%v", response)
}

func runWithDeadline(connection languagepb.LanguageServiceClient, text string) {
	deadlineMs := 2000
	clientDeadline := time.Now().Add(time.Duration(deadlineMs) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), clientDeadline)
	defer cancel()

	request := &languagepb.LanguageDetectorRequest{
		Text: text,
	}
	response, err := connection.DetectWithDeadline(ctx, request)
	if err != nil {
		statusError, ok := status.FromError(err)
		if ok {
			if statusError.Code() == codes.DeadlineExceeded {
				log.Printf("Request has been cancelled becasue it has exceeded deadline %d ms", deadlineMs)
			} else {
				log.Printf("Error detected, type: %v, description: %v", statusError.Code(), statusError.Message())
			}

			return
		}

		log.Println("Error while calling Detect RPC: %v", err)
		return
	}

	log.Printf("Response from Detect Language Service.\n%v", response)
}
