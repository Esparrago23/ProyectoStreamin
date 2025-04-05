package services

import (
    "github.com/IBM/sarama"
    "os"
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
    msg := &sarama.ProducerMessage{
        Topic: os.Getenv("KAFKA_TOPIC"),
        Value: sarama.StringEncoder(string(videoID)),
    }

    _, _, err := ks.producer.SendMessage(msg)
    return err
}