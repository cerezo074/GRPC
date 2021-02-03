package main

import (
	"context"
	"io"
	"math"
	"os"
	"service/rpc/videopb"

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
	videoName := "Tunnel Motions.mp4"
	connection := videopb.NewVideoServiceClient(clientConnection)
	uploadVideo(connection, videoName)
}

func findVideo(filename string, bufferSize int) (*os.File, int64, error) {
	file, err := os.Open("../../assets/client/" + filename)

	if err != nil {
		return nil, 0, err
	}

	stats, err := file.Stat()

	if err != nil {
		return nil, 0, err
	}

	fileSize := stats.Size()

	return file, fileSize, nil
}

func uploadVideo(connection videopb.VideoServiceClient, videoName string) {
	bufferSize := 64 * 1024
	file, fileSize, err := findVideo(videoName, bufferSize)
	if err != nil {
		log.Fatalf("Error opening video file, %v", err)
		return
	}

	requestStream, err := connection.Upload(context.Background())
	if err != nil {
		log.Fatalf("Error connecting RPC Upload method, %v", err)
		return
	}

	if err := runStream(requestStream, videoName, file, bufferSize, fileSize); err != nil {
		log.Fatalf("Error sending data through stream, %v", err)
		return
	}

	if _, err := requestStream.CloseAndRecv(); err != nil {
		log.Fatalf("Error closing RPC Upload method, %v", err)
		return
	}

	log.Printf("Sent %s file!", videoName)
}

func runStream(requestStream videopb.VideoService_UploadClient, videoName string, file *os.File, bufferSize int, fileSize int64) error {
	iterations := int64(math.Ceil(float64(fileSize) / float64(bufferSize)))
	buffer := make([]byte, bufferSize)

	for i := int64(1); i <= iterations; i++ {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		request := &videopb.VideoUploaderRequest{
			Data: &videopb.VideoChunk{
				Content:         buffer[:bytesRead],
				CurrentSecuence: int32(i),
				LastSecuence:    int32(iterations),
				SecuenceSize:    int64(bufferSize),
				Filename:        videoName,
			},
		}

		if err := requestStream.Send(request); err != nil {
			return err
		}

		log.Printf("Sent chunk in secuence %d from total %d secuences", i, iterations)
	}

	return nil
}
