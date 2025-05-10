package examples

import "fmt"

type ChunkMeta struct {
	FileName string
	MD5Hash  string
	Index    int
}

func main() {
	// 1) Using a slice (array-like) you *can* index:
	var slice []ChunkMeta
	slice = append(slice, ChunkMeta{FileName: "file.chunk0", MD5Hash: "abc123", Index: 0})
	slice = append(slice, ChunkMeta{FileName: "file.chunk1", MD5Hash: "def456", Index: 1})

	fmt.Println("slice[0]:", slice[0]) // works, prints the first element
	fmt.Println("slice[1]:", slice[1]) // works, prints the second element

	// 2) Using a channel with a buffer you *cannot* index:
	chunkChan := make(chan ChunkMeta, 2)
	chunkChan <- ChunkMeta{FileName: "file.chunk0", MD5Hash: "abc123", Index: 0}
	chunkChan <- ChunkMeta{FileName: "file.chunk1", MD5Hash: "def456", Index: 1}

	// fmt.Println(chunkChan[0])          // ✗ compile error: cannot index chunkChan
	// fmt.Println(chunkChan[1])          // ✗ compile error: cannot index chunkChan

	// Instead you receive from it in FIFO order:
	first := <-chunkChan
	second := <-chunkChan

	fmt.Println("first from channel:", first)   // the element you sent first
	fmt.Println("second from channel:", second) // the element you sent second
}
