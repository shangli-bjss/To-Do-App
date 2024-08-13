package main

import (
	"fmt"
	"sync"
	"time"
)

func syncDemo() {
	var data int
	var wg sync.WaitGroup
	var mu sync.Mutex

	done := make(chan bool)

	wg.Add(2)

	go func() {
		defer wg.Done()
		for i:=1; i<=9; i+=2 {
			mu.Lock()
			data = i
			fmt.Printf("Odd Goroutine - Setting data to %d\n", data)
			// mu.Unlock()

			time.Sleep(100 * time.Millisecond)

			// mu.Lock()
			fmt.Printf("Odd Goroutine - Current data: %d\n", data)
			mu.Unlock()

			done <- true
		}
	}()

	go func() {
		defer wg.Done()
		for i:=2; i<=10; i+=2 {
			<- done
			
			mu.Lock()
			data = i
			fmt.Printf("Even Goroutine - Setting data to %d\n", data)
			// mu.Unlock()

			time.Sleep(100 * time.Millisecond)

			// mu.Lock()
			fmt.Printf("Even Goroutine - Current data: %d\n", data)
			mu.Unlock()
		}
	}()

	wg.Wait()
	close(done)
}