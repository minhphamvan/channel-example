package main

import (
	"fmt"
	"strings"
	"time"
)

func main() {
	// Khởi tạo 1 publisher
	p := NewPublisher(50, 1*time.Second)

	// Đảm bảo p được đóng trước khi exit
	defer p.Close()

	// `all` subscribe hết tất cả topic
	chanAll := p.Subscribe()

	// Subscribe các topic có "golang"
	chanGolang := p.SubscribeTopic(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "golang")
		}

		return false
	})

	// Print những gì subscriber `all` nhận được
	go func() {
		for msg := range chanAll {
			fmt.Println("==> Chan `all` received:", msg)
		}
	}()

	// Print những gì subscriber `golang` nhận được
	go func() {
		for msg := range chanGolang {
			fmt.Println("==> Chan `golang` received:", msg)
		}
	}()

	for idx := 1; idx <= 10; idx++ {
		time.Sleep(1 * time.Second)

		// Publish ra 2 message
		p.Publish(fmt.Sprintf("Hello world %d", idx))
		p.Publish(fmt.Sprintf("Hello golang %d", idx))
	}

	time.Sleep(10 * time.Second)
}
