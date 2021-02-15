package main

import (
	"context"
	"service/cmd/rpc/languagepb"

	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	fmt.Println("Loading SSL Encryption...")
	certFile := "../cmd/tls/output/ca.crt" // Certificate Authority Trust certificate
	creds, sslErr := credentials.NewClientTLSFromFile(certFile, "")
	if sslErr != nil {
		log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
		return
	}

	opts := grpc.WithTransportCredentials(creds)
	fmt.Println("Running Client...")
	clientConnection, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("Could not connect %v", err)
		return
	}

	defer clientConnection.Close()

	connection := languagepb.NewLanguageServiceClient(clientConnection)
	request := &languagepb.LanguageDetectorRequest{
		// Text: "上海＝南部さやか】世界保健機関（ＷＨＯ）の国際調査団が本格調査を始めた中国湖北省武漢市で、新型コロナウイルスの感染で家族を亡くした遺族に対し、当局からの圧力が強まっている。",
		Text: "Esta variante se propaga más rápido y podría conllevar mayor riesgo de infección, por eso, los investigadores creen que esta nueva variante se propagó de Brasil a Asia.",
	}
	response, err := connection.Detect(context.Background(), request)
	if err != nil {
		log.Fatalf("Error while calling Detect RPC: %v", err)
		return
	}

	log.Printf("Response from Detect Language Service.\n%v", response)
}
