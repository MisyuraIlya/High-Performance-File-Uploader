package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	// Create a sample file with 12 bytes
	sample := []byte("Hello, Gophers")
	if err := os.WriteFile("sample.txt", sample, 0644); err != nil {
		log.Fatal(err)
	}

	f, err := os.Open("sample.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	chunkSize := 5
	buffer := make([]byte, chunkSize)
	chunkIndex := 0

	for {
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}

		// Show the raw buffer and the sliced buffer
		fmt.Printf("Chunk %d: Read %d bytes\n", chunkIndex, n)
		fmt.Printf("  raw buffer:   %q\n", buffer)      // always length 5
		fmt.Printf("  valid slice: %q\n\n", buffer[:n]) // only n bytes

		chunkIndex++
		if err == io.EOF {
			break
		}
	}
}
