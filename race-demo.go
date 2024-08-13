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
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Odd Goroutine - Set data to %d\n", data)
		}
	}()

	go func() {
		defer wg.Done()
		for i:=2; i<10; i+=2{
			data = i
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Even Goroutine - Set data to %d\n", data)
		}
	}()

	wg.Wait()
	pl("Race completed")
}