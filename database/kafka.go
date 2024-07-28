package database

import "github.com/segmentio/kafka-go"

type kafkaKey string
var CtxKey kafkaKey = "kafka"
var KafkaReaderConfig = kafka.ReaderConfig{
    Brokers:  []string{"localhost:9092"},
    Topic:    "messages",
    GroupID:  "consumer1",
    MinBytes: 10e3, // 10KB
    MaxBytes: 10e6, // 10MB
}
var KafkaWriterConfig = kafka.WriterConfig{
    Brokers:  []string{"10.18.0.20:9092"},
    Topic:    "messages",
}
var Reader = kafka.NewReader(KafkaReaderConfig)