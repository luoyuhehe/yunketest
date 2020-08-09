package log

import (
	"github.com/Shopify/sarama"
)

type KafkaLog struct {
	Producer sarama.SyncProducer
	Topic    string
}

// Write
func (lk *KafkaLog) Write(p []byte) (n int, err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = lk.Topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err = lk.Producer.SendMessage(msg)
	if err != nil {
		return
	}
	return
}
