package kafka

import (
    "github.com/segmentio/kafka-go"
    "log"
)

var Writer *kafka.Writer

func Init() {
    Writer = kafka.NewWriter(kafka.WriterConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "example-topic",
    })

    log.Println("Kafka writer created.")
}
