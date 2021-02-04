package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"net"
	"service/rpc/imagepb"

	"github.com/anthonynsimon/bild/effect"
	"google.golang.org/grpc"
)

type ImageHandler struct {
	imagepb.UnsafeImageServiceServer
}

func (handler *ImageHandler) Effect(fileStreamRequest imagepb.ImageService_EffectServer) error {
	err := processStream(fileStreamRequest)
	if err != nil {
		log.Printf("Error receiving chunk data, %v", err)
		return err
	}

	return nil
}

func processStream(fileStreamRequest imagepb.ImageService_EffectServer) error {
	buffer := &bytes.Buffer{}
	currentFileName := ""

	for {
		chunkRequest, err := fileStreamRequest.Recv()
		if err != nil {
			if err == io.EOF {
				err = applyEffect(fileStreamRequest, buffer.Bytes(), currentFileName)
			}
			return err
		}

		log.Printf("Receive new chunk data, secuence %d from total %d", chunkRequest.Data.CurrentSecuence, chunkRequest.Data.LastSecuence)

		if _, err = buffer.Write(chunkRequest.Data.Content); err != nil {
			return err
		}

		if currentFileName == "" {
			currentFileName = chunkRequest.Data.Filename
		} else if currentFileName != chunkRequest.Data.Filename {
			go applyEffect(fileStreamRequest, buffer.Bytes(), currentFileName)
			buffer = &bytes.Buffer{}
			currentFileName = chunkRequest.Data.Filename
		}
	}
}

func applyEffect(responseStream imagepb.ImageService_EffectServer, rawData []byte, imageName string) error {
	img, _, _ := image.Decode(bytes.NewReader(rawData))
	result := effect.EdgeDetection(img, 2.0)
	rawImage, _ := getBytes(result)
	bufferSize := 64 * 1024
	buffer := make([]byte, bufferSize)
	chunks := int64(math.Ceil(float64(len(rawImage)) / float64(bufferSize)))
	dataToSend := bytes.NewReader(rawImage)

	for i := int64(1); i <= chunks; i++ {
		bytesRead, err := dataToSend.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		request := &imagepb.ImageResponse{
			Data: &imagepb.ImageChunk{
				Content:         buffer[:bytesRead],
				CurrentSecuence: int32(i),
				LastSecuence:    int32(chunks),
				SecuenceSize:    int64(bufferSize),
				Filename:        imageName,
			},
		}

		if err := responseStream.Send(request); err != nil {
			return err
		}

		log.Printf("Sent chunk in secuence %d from total %d secuences", i, chunks)
	}

	return nil
}

func getBytes(image image.Image) ([]byte, error) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	if image == nil {
		return nil, fmt.Errorf("image is nil")
	}

	err := jpeg.Encode(w, image, &jpeg.Options{100})

	if err != nil {
		fmt.Println(err)
	}

	return b.Bytes(), nil
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
