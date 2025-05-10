package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"high-performance-file-uploader/internal"
	"log"
	"os"
	"sync"
	"time"
)

const (
	defaultChunkSize = 1024 * 1024
	maxRetries       = 3
)

func loadEnv() error {
	return godotenv.Load()
}

func main() {
	err := loadEnv()
	if err != nil {
		log.Println("No .env file fogo get github.com/joho/godotenv, using default configuration")
	}

	chunkSize := defaultChunkSize
	if size, ok := os.LookupEnv("CHUNK_SIZE"); ok {
		_, err := fmt.Sscanf(size, "%d", &chunkSize)
		if err != nil {
			log.Printf("Error parsing CHUNK_SIZE, using default configuration: %v", err)
		}
	}

	serverURL, ok := os.LookupEnv("SERVER_URL")
	if !ok {
		log.Fatalf("Usage: go run main.go <file_path>")
	}

	filePath := os.Args[1]
	fmt.Println("Loading file ", filePath)
	config := internal.Config{
		ChunkSize: chunkSize,
		ServerURL: serverURL,
	}

	chunker := &internal.DefaultFileChunker{
		ChunkSize: config.ChunkSize,
	}
	uploader := &internal.DefaultUploader{
		ServerURL: config.ServerURL,
	}
	metadataManager := &internal.DefaultMetadataManager{}

	chunks, err := chunker.ChunkFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	metadata, err := metadataManager.LoadMetadata(fmt.Sprintf("%s.metadata.json", filePath))
	if err != nil {
		log.Println("Cloud not load metadata starting fresh.")
		metadata = make(map[string]internal.ChunkMeta)

	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	err = internal.SynchronizeChunks(chunks, metadata, uploader, &wg, &mu)
	if err != nil {
		log.Fatal(err)
	}
	wg.Wait()
	err = metadataManager.SaveMetadata(fmt.Sprintf("%s.metadata.json", filePath), metadata)
	if err != nil {
		log.Fatal(err)
	}

	changeChan := make(chan bool)
	go internal.WatchFile(filePath, changeChan)

	for {
		select {
		case <-changeChan:
			log.Println("File changed, re-chunking and synchonizing....")
			chunks, err = chunker.ChunkFile(filePath)
			if err != nil {
				log.Fatal(err)
			}

			err = internal.SynchronizeChunks(chunks, metadata, uploader, &wg, &mu)
			if err != nil {
				log.Fatal(err)
			}

			wg.Wait()

			err = metadataManager.SaveMetadata(fmt.Sprintf("%s.metadata.json", filePath), metadata)
			if err != nil {
				log.Fatal(err)
			}
		case <-time.After(10 * time.Second):
			log.Println("No changes detected, checking again...")
		}
	}
}
