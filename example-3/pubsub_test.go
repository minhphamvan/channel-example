package main

import (
	"testing"
	"time"
)

func TestNewPublisher(t *testing.T) {
	buffer := 100
	timeout := time.Second
	p := NewPublisher(buffer, timeout)

	if p.buffer != buffer {
		t.Errorf("expected buffer size 100, got %d", p.buffer)
	}

	if p.timeout != timeout {
		t.Errorf("expected timeout %v, got %v", time.Second, p.timeout)
	}

	if len(p.subscribers) != 0 {
		t.Errorf("expected no subscribers, got %d", len(p.subscribers))
	}
}

func TestPublisher_Subcribe_UnSubscribe(t *testing.T) {
	buffer := 100
	timeout := time.Second
	p := NewPublisher(buffer, timeout)

	sub := p.Subscribe()
	if len(p.subscribers) != 1 {
		t.Errorf("expected 1 subscriber, got %d", len(p.subscribers))
	}

	p.UnSubscribe(sub)
	if len(p.subscribers) != 0 {
		t.Errorf("expected 0 subscribers, got %d", len(p.subscribers))
	}
}

func TestPublisher_Publish(t *testing.T) {
	buffer := 100
	timeout := time.Second
	p := NewPublisher(buffer, timeout)

	sub := p.Subscribe()

	go p.Publish("test message")

	select {
	case msg := <-sub:
		if msg != "test message" {
			t.Errorf("expected 'test message', got %v", msg)
		}
	case <-time.After(time.Second):
		t.Error("did not receive message in time")
	}
}

func TestPublisher_Close(t *testing.T) {
	buffer := 100
	timeout := time.Second
	p := NewPublisher(buffer, timeout)

	sub := p.Subscribe()

	p.Close()

	if len(p.subscribers) != 0 {
		t.Errorf("expected 0 subscribers after close, got %d", len(p.subscribers))
	}

	select {
	case <-sub:
		// expected channel to be closed
	default:
		t.Error("expected subscriber channel to be closed")
	}
}
