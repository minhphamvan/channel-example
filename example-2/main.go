package main

import (
	"fmt"
	"sync"
	"time"
)

func startProducers(x int, pub *Publisher, wg *sync.WaitGroup) {
	for i := 1; i <= x; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for j := 0; j < 5; j++ {
				value := fmt.Sprintf("Producer %d, message %d", id, j)
				fmt.Printf("Producer %d publishing: %s\n", id, value)

				pub.Publish(value)
				time.Sleep(time.Millisecond * 500)
			}
		}(i)
	}
}

func startConsumers(y int, pub *Publisher, wg *sync.WaitGroup) {
	for i := 1; i <= y; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			sub := pub.Subscribe()
			for value := range sub {
				fmt.Printf("Consumer %d received: %v\n", id, value)
				time.Sleep(time.Millisecond * 500)
			}
		}(i)
	}
}

func main() {
	var x, y int
	fmt.Print("Enter number of producers:")
	fmt.Scan(&x)
	fmt.Print("Enter number of consumers:")
	fmt.Scan(&y)

	// Create a new Publisher with a buffer size of 10 and a timeout of 100 milliseconds
	pub := NewPublisher(10, 100*time.Millisecond)

	// WaitGroup to wait for all goroutines
	var wg sync.WaitGroup

	// Start producers (x)
	startProducers(x, pub, &wg)

	// Start consumers (y)
	startConsumers(y, pub, &wg)

	// Wait for all producers to finish
	wg.Wait()

	// Close the publisher and all subscriber channels
	pub.Close()

	fmt.Println("All producers and consumers have finished.")
}

// func main() {
// 	// Khởi tạo 1 publisher
// 	p := NewPublisher(50, 1*time.Second)
//
// 	// Đảm bảo p được đóng trước khi exit
// 	defer p.Close()
//
// 	// `all` subscribe hết tất cả topic
// 	chanAll := p.Subscribe()
//
// 	// Subscribe các topic có "golang"
// 	chanGolang := p.SubscribeTopic(func(v interface{}) bool {
// 		if s, ok := v.(string); ok {
// 			return strings.Contains(s, "golang")
// 		}
//
// 		return false
// 	})
//
// 	// Print những gì subscriber `all` nhận được
// 	go func() {
// 		for msg := range chanAll {
// 			fmt.Println("==> Chan `all` received:", msg)
// 		}
// 	}()
//
// 	// Print những gì subscriber `golang` nhận được
// 	go func() {
// 		for msg := range chanGolang {
// 			fmt.Println("==> Chan `golang` received:", msg)
// 		}
// 	}()
//
// 	for idx := 1; idx <= 10; idx++ {
// 		time.Sleep(1 * time.Second)
//
// 		// Publish ra 2 message
// 		p.Publish(fmt.Sprintf("Hello world %d", idx))
// 		p.Publish(fmt.Sprintf("Hello golang %d", idx))
// 	}
//
// 	time.Sleep(10 * time.Second)
// }
