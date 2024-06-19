package main

import (
	"fmt"
	"time"
)

func producer(ch chan int) {
	for i := 1; i <= 10; i++ {
		time.Sleep(time.Second)

		fmt.Println("Producer: sending", i)
		ch <- i
	}
	close(ch)
}

func consumer(ch <-chan int) {
	for i := range ch {
		fmt.Println("===> Consumer: receiving", i)
		fmt.Println()
	}
}

func main() {
	ch := make(chan int, 3)
	go producer(ch)
	go consumer(ch)
	time.Sleep(time.Second * 20)
}
