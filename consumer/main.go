package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("Failed to connect to kafka consumer: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Printf("Error closing consumer: %v", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("orders", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Failed to consume partition: %v", err)
	}
	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Printf("Error closing partition consumer: %v", err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
	log.Println("Consumer started. Waiting for messages...")

ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message from partition %d at offset %d", msg.Partition, msg.Offset)
			log.Printf("Raw message: %s", string(msg.Value))
			consumed++
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming message: %v", err)
		case <-signals:
			log.Println("Interrupt received, shutting down...")
			break ConsumerLoop
		}
	}

	log.Printf("Total consumed: %d messages", consumed)
}
