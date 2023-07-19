package producer

import (
	"order-and-pay/iface"

	"github.com/IBM/sarama"
)

type KafkaProducer struct {
	sarama.SyncProducer
	logger iface.Ilogger
}

func NewKafkaProducer(brokersUrl []string, l iface.Ilogger) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	syncProducer, err := sarama.NewSyncProducer(brokersUrl, config)
	if err != nil {
		return nil, err
	}
	return &KafkaProducer{syncProducer, l}, nil
}

func (p *KafkaProducer) PushMessageToQueue(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := p.SendMessage(msg)
	if err != nil {
		return err
	}
	p.logger.Infof("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}
