package main

import (
	"fmt"
	"time"
)

func Producer(ch chan int, done chan struct{}) {
	for i := 1; i <= 10; i++ {
		time.Sleep(time.Second)

		fmt.Println("Producer: sending", i)
		ch <- i
	}

	close(ch)
	done <- struct{}{}
}

func Consumer(ch <-chan int, done chan struct{}) {
	for i := range ch {
		fmt.Println("===> Consumer: receiving", i)
		fmt.Println()
	}

	done <- struct{}{}
}

func main() {
	ch := make(chan int, 3)
	done := make(chan struct{})

	go Producer(ch, done)
	go Consumer(ch, done)

	<-done
	<-done
}
