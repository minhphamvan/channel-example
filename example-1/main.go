package main

import (
	"fmt"
	"time"
)

func producer(ch chan int, done chan struct{}) {
	for i := 1; i <= 10; i++ {
		time.Sleep(time.Second)

		fmt.Println("Producer: sending", i)
		ch <- i
	}

	close(ch)
	done <- struct{}{}
}

func consumer(ch <-chan int, done chan struct{}) {
	for i := range ch {
		fmt.Println("===> Consumer: receiving", i)
		fmt.Println()
	}

	done <- struct{}{}
}

func main() {
	ch := make(chan int, 3)
	done := make(chan struct{})

	go producer(ch, done)
	go consumer(ch, done)

	<-done
	<-done
}
