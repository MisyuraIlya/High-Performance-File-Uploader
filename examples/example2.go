package examples

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex

	// Goroutine A
	go func() {
		mu.Lock()
		fmt.Println("A: got lock, working for 2s…")
		time.Sleep(2 * time.Second)
		fmt.Println("A: releasing lock")
		mu.Unlock()
	}()

	// Give A a moment to grab the lock
	time.Sleep(100 * time.Millisecond)

	// Goroutine B
	go func() {
		fmt.Println("B: trying to get lock")
		mu.Lock() // <— blocks here until A calls Unlock()
		fmt.Println("B: got lock!")
		mu.Unlock()
	}()

	// Goroutine C
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("C: tick %d\n", i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Let everything run
	time.Sleep(4 * time.Second)
}

//A: got lock, working for 2
//C: tick 0
//C: tick 1
//B: trying to get lock
//C: tick 2
//A: releasing lock
//B: got lock!
//C: tick 3
//C: tick 4
