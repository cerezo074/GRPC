package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"os"
	"service/rpc/imagepb"

	"google.golang.org/grpc"
)

type ImageHandler struct {
	imagepb.UnsafeImageServiceServer
}

func (*ImageHandler) findImage(filename string, bufferSize int) (*os.File, int64, error) {
	file, err := os.Open("../../assets/server/" + filename)

	if err != nil {
		log.Fatalf("Error opening image file: %v", err)
		return nil, 0, err
	}

	stats, err := file.Stat()

	if err != nil {
		log.Fatalf("Error opening image file: %v", err)
		return nil, 0, err
	}

	fileSize := stats.Size()

	return file, fileSize, nil
}

func (handler *ImageHandler) Download(request *imagepb.ImageDownloaderRequest, responseStream imagepb.ImageService_DownloadServer) error {
	fmt.Printf("Language detector function was invoked with %v", request)
	imageName := request.GetImageName()
	bufferSize := 64 * 1024
	file, fileSize, err := handler.findImage(imageName, bufferSize)
	if err != nil {
		log.Fatalf("Error opening image file, %v", err)
		return err
	}

	iterations := int64(math.Ceil(float64(fileSize) / float64(bufferSize)))
	buffer := make([]byte, bufferSize)

	for i := int64(1); i <= iterations; i++ {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			log.Println("Finish sending file:" + imageName)
			return nil
		} else if err != nil {
			log.Fatalf("Error reading image file, %v", err)
			return err
		}

		response := &imagepb.ImageDownloaderResponse{
			Result: &imagepb.ImageChunk{
				Content:         buffer[:bytesRead],
				CurrentSecuence: int32(i),
				LastSecuence:    int32(iterations),
				SecuenceSize:    int64(bufferSize),
			},
		}

		if err := responseStream.Send(response); err != nil {
			log.Fatalf("Error sending response, %v", err)
			return err
		}

		log.Printf("Sent chunk in secuence %d from total %d secuences", i, iterations)
	}

	return nil
}

func main() {
	listener, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		log.Fatalf("Failed to listen, %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	imagepb.RegisterImageServiceServer(grpcServer, &ImageHandler{})
	fmt.Println("Server Listening...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve, %v", err)
	}
}
