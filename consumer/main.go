package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Kafka consumer configuration
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	topic := os.Getenv("KAFKA_TOPIC")

	// Create Kafka consumer
	consumer, err := sarama.NewConsumer([]string{kafkaBrokers}, nil)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Create partition consumer
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("Error creating partition consumer: %v", err)
	}
	defer partitionConsumer.Close()

	// Start HTTP server
	go func() {
		http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Consumer is running")
		})
		log.Printf("Consumer service starting on port 8081...")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start consuming messages
	log.Printf("Consuming messages from Kafka topic: %s\n", topic)
	for msg := range partitionConsumer.Messages() {
		log.Printf("Received message: %s\n", string(msg.Value))
		videoID := string(msg.Value)
		if videoID == "" {
			log.Printf("Received empty videoID, skipping processing")
			continue
		}
		processVideo(videoID)
	}
}

func processVideo(videoID string) {
	log.Printf("Processing video with ID: %s\n", videoID)
	// Add video processing logic here
}
