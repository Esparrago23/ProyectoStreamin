package services

import (
    "github.com/IBM/sarama"
    "os"
    "log"
    "strconv"
)

type KafkaService struct {
    producer sarama.SyncProducer
}

func NewKafkaService() (*KafkaService, error) {
    config := sarama.NewConfig()
    config.Producer.Return.Successes = true
    
    producer, err := sarama.NewSyncProducer([]string{os.Getenv("KAFKA_BROKERS")}, config)
    if err != nil {
        return nil, err
    }

    return &KafkaService{producer: producer}, nil
}

func (ks *KafkaService) SendVideoForProcessing(videoID uint) error {
    // Convert the videoID to string for sending it as a message
    videoIDString := strconv.Itoa(int(videoID))

    // Log the video ID being sent
    log.Printf("Sending videoID %s to Kafka topic %s", videoIDString, os.Getenv("KAFKA_TOPIC"))
    
    // Create a message with the videoID to be sent to Kafka
    msg := &sarama.ProducerMessage{
        Topic: os.Getenv("KAFKA_TOPIC"),
        Value: sarama.StringEncoder(videoIDString),
    }

    // Send the message
    partition, offset, err := ks.producer.SendMessage(msg)
    if err != nil {
        log.Printf("Failed to send message to Kafka: %v", err)
        return err
    }

    // Log the success and show the partition and offset
    log.Printf("Message sent successfully to Kafka topic %s, partition %d, offset %d", os.Getenv("KAFKA_TOPIC"), partition, offset)
    return nil
}
