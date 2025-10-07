package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Order struct {
	OrderID  string `json:"order_id"`
	User     string `json:"user"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

func (o Order) Encode() ([]byte, error) {
	return json.Marshal(o)
}

var (
	users = []string{"bughowi", "alice", "bob", "charlie", "diana"}
	items = []string{
		"mie gacoan", "mushroom pizza", "nasi goreng", "chicken katsu",
		"ramen bowl", "beef burger", "pad thai", "sushi roll",
	}
)

func randomOrder() Order {
	id, _ := uuid.NewV7()
	return Order{
		OrderID:  id.String(),
		User:     users[rand.Intn(len(users))],
		Item:     items[rand.Intn(len(items))],
		Quantity: rand.Intn(10) + 1,
	}
}

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewAsyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Kafka: %v", err)
	}
	defer func() {
		log.Println("Closing producer...")
		producer.AsyncClose()
		log.Println("Producer closed.")
	}()

	// Graceful shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)

	// Listen for delivery reports
	go func() {
		for {
			select {
			case s := <-producer.Successes():
				if order, ok := s.Metadata.(Order); ok {
					log.Printf("✅ Sent order %s → partition %d, offset %d",
						order.OrderID, s.Partition, s.Offset)
				}
			case e := <-producer.Errors():
				log.Printf("❌ Failed to send message: %v", e.Err)
			}
		}
	}()

	log.Println("🚀 Producer started. Sending random orders every 5 seconds...")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	sendOrder(producer)

	for {
		select {
		case <-ticker.C:
			sendOrder(producer)
		case <-sigchan:
			log.Println("👋 Interrupt received, exiting...")
			return
		}
	}
}

func sendOrder(producer sarama.AsyncProducer) {
	order := randomOrder()
	value, _ := order.Encode()

	msg := &sarama.ProducerMessage{
		Topic:    "orders",
		Value:    sarama.ByteEncoder(value),
		Metadata: order,
	}

	log.Printf("→ Sending order %s: %s (qty: %d) for %s",
		order.OrderID, order.Item, order.Quantity, order.User)

	producer.Input() <- msg
}
