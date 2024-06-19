package main

import (
	"fmt"
	"sync"
	"time"
)

type (
	subscriber chan interface{}
	topicFunc  func(v interface{}) bool
)

type Publisher struct {
	mu          sync.RWMutex             // ReadWrite Mutex
	buffer      int                      // kích thước hàng đợi
	timeout     time.Duration            // timeout cho việc publishing
	subscribers map[subscriber]topicFunc // subscriber đã subscribe vào topic nào
}

func NewPublisher(buffer int, timeout time.Duration) *Publisher {
	return &Publisher{
		mu:          sync.RWMutex{},
		buffer:      buffer,
		timeout:     timeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// Thêm subscriber mới, subscribe các topic đã được filter lọc
func (p *Publisher) SubscribeTopic(topicFunc topicFunc) chan interface{} {
	ch := make(chan interface{}, p.buffer)

	p.mu.Lock()
	p.subscribers[ch] = topicFunc
	p.mu.Unlock()

	return ch
}

// Thêm subscriber mới, đăng ký hết tất cả topic
func (p *Publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

// Hủy subscribe
func (p *Publisher) UnSubscribe(sub chan interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	delete(p.subscribers, sub)
	close(sub)
}

// Publish ra 1 topic
func (p *Publisher) Publish(v interface{}) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// Gửi 1 topic có thể duy trì trong thời gian chờ wg
func (p *Publisher) sendTopic(sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
		fmt.Printf("Send value %v to subscriber %v\n\n", v, sub)
	case <-time.After(p.timeout):
	}
}

// Close publisher và close tất cả các subscriber
func (p *Publisher) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}
