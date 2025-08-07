package main

import (
	"f1term/internal/fetchData"
	"fmt"
	"time"
)

func main() {
	done := make(chan struct{})

	go fetchData.FetchByYears(done)

	for {
		select {
		case <-done:
			fmt.Println("Terminé.")
			return
		default:
			fmt.Println("Waiting...")
			time.Sleep(1 * time.Second)
		}
	}
}
