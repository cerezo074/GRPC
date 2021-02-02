package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"service/rpc/videopb"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type VideoHandler struct {
	videopb.UnsafeVideoServiceServer
}

func (handler *VideoHandler) Upload(fileStreamRequest videopb.VideoService_UploadServer) error {
	buffer, filename, err := processStream(fileStreamRequest)
	if err != nil {
		log.Printf("Error receiving chunk data, %v", err)
		return err
	}

	err = processBuffer(buffer, filename)
	if err != nil {
		log.Printf("Error creating file, %v", err)
		return err
	}

	fileStreamRequest.SendAndClose(&emptypb.Empty{})
	return nil
}

func processStream(fileStreamRequest videopb.VideoService_UploadServer) (*bytes.Buffer, string, error) {
	buffer := &bytes.Buffer{}
	filename := ""

	for {
		chunkRequest, err := fileStreamRequest.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, "", err
		}

		log.Printf("Receive new chunk data, secuence %d from total %d", chunkRequest.Data.CurrentSecuence, chunkRequest.Data.LastSecuence)

		if _, err = buffer.Write(chunkRequest.Data.Content); err != nil {
			return nil, "", err
		}

		if filename == "" {
			filename = chunkRequest.Data.Filename
		}
	}

	return buffer, filename, nil
}

func processBuffer(buffer *bytes.Buffer, videoName string) error {
	basePath := "../../assets/server/"
	videoFile := basePath + videoName
	os.Remove(videoFile)

	if err := ioutil.WriteFile(videoFile, buffer.Bytes(), 777); err != nil {
		return err
	}

	log.Println("File " + videoName + " created successfully")
	return nil
}

func main() {
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen, %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	videopb.RegisterVideoServiceServer(grpcServer, &VideoHandler{})
	fmt.Println("Server Listening...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve, %v", err)
	}
}
