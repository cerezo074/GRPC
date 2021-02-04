package main

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"math"
	"os"
	"service/rpc/imagepb"

	"fmt"
	"log"

	"google.golang.org/grpc"
)

type uploadFile struct {
	name   string
	data   *os.File
	size   int64
	chunks int64
}

func main() {
	fmt.Println("Running Client...")

	clientConnection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
		return
	}

	defer clientConnection.Close()
	videoNames := []string{"zoro.jpg", "grimmjow.jpg"}
	connection := imagepb.NewImageServiceClient(clientConnection)
	uploadImages(connection, videoNames...)
}

func uploadImages(connection imagepb.ImageServiceClient, imageNames ...string) {
	bufferSize := 64 * 1024
	files := newUploadFiles(bufferSize, imageNames...)
	if len(files) <= 0 {
		return
	}

	uploadStream, err := connection.Effect(context.Background())
	if err != nil {
		log.Fatalf("Error connecting RPC Upload method, %v", err)
		return
	}

	waitChannel := make(chan struct{})

	go uploadImageFiles(uploadStream, files, bufferSize)
	go receiveImageFiles(waitChannel, uploadStream)

	<-waitChannel
}

func newUploadFiles(bufferSize int, imageNames ...string) []uploadFile {
	files := make([]uploadFile, len(imageNames))
	for _, imageName := range imageNames {
		file, size, err := findVideo(imageName)
		if err != nil {
			log.Fatalf("Error opening file, %v", err)
			continue
		}

		chunks := int64(math.Ceil(float64(size) / float64(bufferSize)))
		uploadFile := uploadFile{name: imageName, data: file, size: size, chunks: chunks}
		files = append(files, uploadFile)
	}

	return files
}

func findVideo(filename string) (*os.File, int64, error) {
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

func uploadImageFiles(uploadStream imagepb.ImageService_EffectClient, files []uploadFile, bufferSize int) {
	for _, file := range files {
		if err := runStream(uploadStream, file, bufferSize); err != nil {
			log.Fatalf("Upload %s file has been canceled due to an error %v, ", file.name, err)
			uploadStream.CloseSend()
			return
		}
	}

	uploadStream.CloseSend()
	log.Println("Uploaded all files")
}

func runStream(requestStream imagepb.ImageService_EffectClient, file uploadFile, bufferSize int) error {
	buffer := make([]byte, bufferSize)

	for i := int64(1); i <= file.chunks; i++ {
		bytesRead, err := file.data.Read(buffer)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		request := &imagepb.ImageRequest{
			Data: &imagepb.ImageChunk{
				Content:         buffer[:bytesRead],
				CurrentSecuence: int32(i),
				LastSecuence:    int32(file.chunks),
				SecuenceSize:    int64(bufferSize),
				Filename:        file.name,
			},
		}

		if err := requestStream.Send(request); err != nil {
			return err
		}

		log.Printf("Sent chunk in secuence %d from total %d secuences", i, file.chunks)
	}

	return nil
}

func receiveImageFiles(waitChannel chan<- struct{}, uploadStream imagepb.ImageService_EffectClient) {
	buffer := &bytes.Buffer{}
	currentFileName := ""

	for {
		chunkResponse, err := uploadStream.Recv()
		if err == io.EOF {
			saveFile(buffer.Bytes(), currentFileName)
			log.Println("Closing connection with RPC, stream has sent all data")
			break
		} else if err != nil {
			log.Printf("Error while receive chunk data, %v", err)
			break
		}

		log.Printf("Receive new chunk data, secuence %d from total %d", chunkResponse.Data.CurrentSecuence, chunkResponse.Data.LastSecuence)

		if _, err = buffer.Write(chunkResponse.Data.Content); err != nil {
			log.Printf("Error writing chunk data, %v", err)
			break
		}

		if currentFileName == "" {
			currentFileName = chunkResponse.Data.Filename
		} else if currentFileName != chunkResponse.Data.Filename {
			go saveFile(buffer.Bytes(), currentFileName)
			buffer = &bytes.Buffer{}
			currentFileName = chunkResponse.Data.Filename
		}
	}

	close(waitChannel)
}

func saveFile(buffer []byte, imageName string) {
	basePath := "../../assets/result/"
	imageFile := basePath + imageName
	os.Remove(imageFile)

	if err := ioutil.WriteFile(imageFile, buffer, 777); err != nil {
		log.Printf("Error writing file, %v", err)
		return
	}

	log.Println("File " + imageName + " created successfully inside images folder")
}
