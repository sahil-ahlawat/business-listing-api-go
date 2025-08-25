package main

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal("❌ NATS connection failed:", err)
	}
	defer nc.Drain()

	// Subscribe
	_, err = nc.Subscribe("test.subject", func(m *nats.Msg) {
		fmt.Println("✅ Received message:", string(m.Data))
	})
	if err != nil {
		log.Fatal("❌ Failed to subscribe:", err)
	}

	// Give it a second
	time.Sleep(1 * time.Second)

	// Publish
	err = nc.Publish("test.subject", []byte("Hello from Go!"))
	if err != nil {
		log.Fatal("❌ Failed to publish:", err)
	}

	// Wait to receive
	time.Sleep(1 * time.Second)
}
