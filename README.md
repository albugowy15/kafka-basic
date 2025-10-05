# Kafka Producer & Consumer in Go

A simple implementation of Apache Kafka producer and consumer using Go and the Sarama library. This project demonstrates basic message publishing and consumption patterns with randomly generated order data.

## Features

- **Async Producer**: Sends random order messages to Kafka every 5 seconds
- **Consumer**: Consumes and displays order messages from Kafka
- **Automatic Message Generation**: Creates random orders with varied users, items, and quantities
- **Graceful Shutdown**: Handles interrupt signals (Ctrl+C) cleanly

## Prerequisites

Before running this project, ensure you have the following installed:

- **Docker** - [Install Docker](https://docs.docker.com/get-docker/)
- **Docker Compose** - [Install Docker Compose](https://docs.docker.com/compose/install/)
- **Go 1.24** - [Install Go](https://golang.org/doc/install)

## Project Structure

```
.
├── docker-compose.yml
├── producer/
│   └── main.go
├── consumer/
│   └── main.go
└── README.md
```

## Getting Started

### 1. Start Kafka with Docker

Run Kafka container using Docker Compose:

```bash
docker compose up -d
```

This will start Kafka broker on port `9092`

To verify the containers are running:

```bash
docker compose ps
```

### 2. Run the Consumer

Open a terminal and navigate to the consumer directory:

```bash
cd consumer
go run main.go
```

The consumer will start listening for messages on the `orders` topic. You should see:

```
2025/10/06 01:02:56 Consumer started. Waiting for messages...
2025/10/06 01:02:59 Consumed message from partition 0 at offset 351346
2025/10/06 01:02:59 Raw message: {"order_id":"f9a2212c-0c65-f6b5-5973-d0bcaa382188","user":"alice","item":"beef burger","quantity":9}
2025/10/06 01:03:04 Consumed message from partition 0 at offset 351347
2025/10/06 01:03:04 Raw message: {"order_id":"4e2e513f-6014-06ff-522c-9602f5d196cc","user":"diana","item":"ramen bowl","quantity":4}
```

### 3. Run the Producer

Open another terminal and navigate to the producer directory:

```bash
cd producer
go run main.go
```

The producer will start sending random order messages every 5 seconds. You should see output like:

```
2025/10/06 01:02:49 🚀 Producer started. Sending random orders every 5 seconds...
2025/10/06 01:02:49 → Sending order 206c33be-3013-7fdd-1f58-ab30476e142e: mushroom pizza (qty: 1) for charlie
2025/10/06 01:02:49 ✅ Sent order 206c33be-3013-7fdd-1f58-ab30476e142e → partition 0, offset 351344
2025/10/06 01:02:54 → Sending order a6a49da5-feca-5339-2267-78b5cd349682: chicken katsu (qty: 3) for diana
2025/10/06 01:02:54 ✅ Sent order a6a49da5-feca-5339-2267-78b5cd349682 → partition 0, offset 351345
```

### 4. Stop the Applications

To stop the producer or consumer, press `Ctrl+C` in their respective terminals.

To stop Kafka

```bash
docker compose down
```

## Order Message Format

Messages are JSON-encoded with the following structure:

```json
{
  "order_id": "uuid-string",
  "user": "username",
  "item": "item-name",
  "quantity": 1-10
}
```
