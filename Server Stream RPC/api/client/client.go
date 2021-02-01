package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"os"
	"service/rpc/imagepb"

	"fmt"
	"log"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Running Client...")

	clientConnection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
		return
	}

	defer clientConnection.Close()
	imageName := "lake.jpg"
	connection := imagepb.NewImageServiceClient(clientConnection)
	buffer, err := downloadImage(connection, imageName)

	if buffer != nil {
		process(buffer, imageName)
	}
}

func downloadImage(connection imagepb.ImageServiceClient, imageName string) (*bytes.Buffer, error) {
	request := imagepb.ImageDownloaderRequest{
		ImageName: imageName,
	}

	fileStreamResponse, err := connection.Download(context.Background(), &request)
	if err != nil {
		log.Fatalf("Error while calling Download RPC, %v", err)
		return nil, err
	}

	buffer := &bytes.Buffer{}

	for {
		chunkResponse, err := fileStreamResponse.Recv()
		if err == io.EOF {
			log.Println("Closing connection with RPC, stream has sent all data")
			break
		} else if err != nil {
			log.Printf("Error while receive chunk data, %v", err)
			return nil, err
		}

		log.Printf("Receive new chunk data, secuence %d from total %d", chunkResponse.Result.CurrentSecuence, chunkResponse.Result.LastSecuence)

		if _, err = buffer.Write(chunkResponse.Result.Content); err != nil {
			log.Printf("Error writing chunk data, %v", err)
			return nil, err
		}
	}

	return buffer, nil
}

func process(buffer *bytes.Buffer, imageName string) {
	imageFile := "../../assets/client/" + imageName
	os.Remove(imageFile)

	if err := ioutil.WriteFile(imageFile, buffer.Bytes(), 777); err != nil {
		log.Printf("Error writing file, %v", err)
		return
	}

	log.Println("File " + imageName + " created successfully inside images folder")
}
