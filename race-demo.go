package main

import (
	"fmt"
	"sync"
	"time"
)

func raceDemo() {
	var data int
	var wg sync.WaitGroup

	wg.Add(2)

	go func(){
		defer wg.Done()
		for i:=1; i<9; i+=2 {
			data = i
			fmt.Printf("Odd Goroutine - Setting data to %d...\n", data)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Odd Goroutine - Set data to %d\n", data)
		}
	}()

	go func() {
		defer wg.Done()
		for i:=2; i<10; i+=2{
			data = i
			fmt.Printf("Even Goroutine - Setting data to %d...\n", data)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Even Goroutine - Set data to %d\n", data)
		}
	}()

	wg.Wait()
	fmt.Println("Race completed")
}