package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// producer: liên tục tạo ra một chuỗi số nguyên dựa trên bội số factor và đưa vào channel
func Producer(factor int, out chan<- int) {
	for i := 1; i <= 10; i++ {
		time.Sleep(time.Second * 1)

		val := i * factor
		out <- val

		fmt.Printf("Produced %v from factor %v\n", val, factor)
	}
}

// consumer: liên tục lấy các số từ channel ra để print
func Consumer(in <-chan int) {
	for v := range in {
		fmt.Println("==> Received:", v)
	}
}

func main() {
	ch := make(chan int, 64)

	// // tạo một chuỗi số với bội số 2
	go Producer(2, ch)

	// tạo một chuỗi số với bội số 3
	go Producer(3, ch)

	// tạo consumer
	go Consumer(ch)

	// Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}
