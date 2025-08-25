// File: fitness/nats/nats_test.go
package nats_test

import (
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

func TestNATSStream(t *testing.T) {
	// Connect to local NATS server
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Drain()

	// Create JetStream context
	js, err := nc.JetStream()
	if err != nil {
		t.Fatalf("Failed to get JetStream context: %v", err)
	}

	streamName := "TEST_STREAM"
	subject := "test.subject"
	message := "Hello from independent test"

	// Add stream
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.MemoryStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		t.Fatalf("Stream creation failed: %v", err)
	}

	// Subscribe and wait for message
	done := make(chan bool)

	_, err = js.Subscribe(subject, func(msg *nats.Msg) {
		if string(msg.Data) != message {
			t.Errorf("Expected %q, got %q", message, msg.Data)
		}
		msg.Ack()
		done <- true
	}, nats.ManualAck())

	if err != nil {
		t.Fatalf("Subscription failed: %v", err)
	}

	// Publish message
	_, err = js.Publish(subject, []byte(message))
	if err != nil {
		t.Fatalf("Publish failed: %v", err)
	}

	// Wait for message or timeout
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Message not received in time")
	}
}
