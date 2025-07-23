package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/haoran-mc/golib/pkg/env"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var (
	brokerAddress = "127.0.0.1:9092"
	topic         = "api-assets"
	groupID       = "api-asset-processor"
	username      = "admin"
	password      = env.GetEnv("SWPASSWORD", "")
)

var messages = []string{
	"one", "two", "three", "four", "five",
	"six", "seven", "eight", "nine", "ten",
}

func createDialer() *kafka.Dialer {
	mechanism, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		log.Fatalf("Failed to create SCRAM mechanism: %v", err)
	}

	return &kafka.Dialer{
		Timeout:       10 * time.Second,
		SASLMechanism: mechanism,
		TLS:           nil, // 如果 Kafka 启用了 TLS，否则设为 nil
	}
}

func produce() {
	dialer := createDialer()
	writer := kafka.Writer{
		Addr:         kafka.TCP(brokerAddress),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Transport: &kafka.Transport{
			SASL: dialer.SASLMechanism,
			TLS:  dialer.TLS,
		},
	}

	defer writer.Close()

	fmt.Println("Producing messages...")
	for _, msg := range messages {
		err := writer.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte("key"),
			Value: []byte(msg),
		})
		if err != nil {
			log.Fatalf("Failed to write message: %v", err)
		}
	}
	fmt.Println("Produced 10 messages")
}

func consume() {
	dialer := createDialer()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{brokerAddress},
		Topic:           topic,
		GroupID:         groupID,
		MinBytes:        1,
		MaxBytes:        10e6,
		CommitInterval:  time.Second, // 自动提交偏移量间隔
		Dialer:          dialer,
		StartOffset:     kafka.FirstOffset,
		ReadLagInterval: -1, // 禁用 lag 监控
	})

	defer reader.Close()

	fmt.Println("Consuming messages...")
	count := 0
	for count < len(messages) {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Failed to read message: %v", err)
		}
		fmt.Printf("Message received: %s\n", string(m.Value))
		count++
	}
}

func main() {
	createTopic()
	// 先写入，再读取
	produce()
	time.Sleep(2 * time.Second)
	// consume()
	consumeWithoutGroup()
}

func createTopic() {
	dialer := createDialer()
	conn, err := dialer.DialContext(context.Background(), "tcp", brokerAddress)
	if err != nil {
		log.Fatalf("Failed to dial Kafka: %v", err)
	}
	defer conn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}

	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatalf("Failed to create topic: %v", err)
	}
	fmt.Println("Topic created successfully")
}

func consumeWithoutGroup() {
	dialer := createDialer()
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       topic,
		Partition:   0,
		MinBytes:    1,
		MaxBytes:    10e6,
		StartOffset: kafka.FirstOffset,
		Dialer:      dialer,
	})

	defer reader.Close()

	fmt.Println("Consuming messages (no group)...")
	count := 0
	for count < len(messages) {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("Failed to read message: %v", err)
		}
		fmt.Printf("Message received: %s\n", string(m.Value))
		count++
	}
}
