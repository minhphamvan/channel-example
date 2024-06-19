package main

import (
	"testing"
	"time"
)

func Test_Producer(t *testing.T) {
	ch := make(chan int, 10)
	done := make(chan struct{})

	go Producer(ch, done)

	select {
	case <-done:
		close(done)
	case <-time.After(12 * time.Second):
		t.Error("producer timed out")
	}

	// Validate
	expected := 10
	count := 0
	for range ch {
		count++
	}

	if count != expected {
		t.Errorf("expected %d items produced, got %d", expected, count)
	}
}

func Test_Consumer(t *testing.T) {
	ch := make(chan int, 10)
	done := make(chan struct{})

	go Consumer(ch, done)

	for i := 1; i <= 10; i++ {
		time.Sleep(time.Second)
		ch <- i
	}
	close(ch)

	select {
	case <-done:
		close(done)
	case <-time.After(12 * time.Second):
		t.Error("consumer timed out")
	}
}
