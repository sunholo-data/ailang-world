package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	shared := 0
	for range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10_000 {
				shared++ // Deliberate: known-positive control for the race detector.
			}
		}()
	}
	wg.Wait()
	fmt.Println("control finished, shared =", shared)
}
