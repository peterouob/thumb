package kafka

import (
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	prom "thumb/pkg/prometheus"
	"time"
)

var Pd *Producer

type Producer struct {
	sarama.AsyncProducer
}

func NewProducer() {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForLocal
	conf.Producer.Return.Errors = true
	conf.Producer.Compression = sarama.CompressionZSTD
	conf.Producer.Flush.Messages = 10
	conf.Producer.Flush.Frequency = 10 * time.Millisecond
	conf.Producer.Flush.Bytes = 1024
	conf.Producer.Retry.Backoff = 500 * time.Millisecond
	producer, err := sarama.NewAsyncProducer([]string{"192.168.0.100:9092"}, conf)
	if errors.Is(err, sarama.ErrClosedClient) {
		panic(errors.New("error in NewAsyncProducer :" + err.Error()))
	}
	pd := &Producer{AsyncProducer: producer}
	pd.handleResponses()
	Pd = pd
}

func (p *Producer) handleResponses() {
	go func() {
		for {
			select {
			case err := <-p.AsyncProducer.Errors():
				prom.KafkaRequestTotal.WithLabelValues(err.Msg.Topic, "failed", err.Err.Error()).Inc()
			case success := <-p.AsyncProducer.Successes():
				msg, err := success.Value.Encode()
				if err != nil {
					prom.KafkaRequestTotal.WithLabelValues(success.Topic, "failed", err.Error()).Inc()
				} else {
					prom.KafkaRequestTotal.WithLabelValues(success.Topic, "success", string(msg)).Inc()
				}
			}
		}
	}()
}

func (p *Producer) Send(topic string, value any) {
	v, err := json.Marshal(value)
	if err != nil {
		prom.KafkaRequestTotal.WithLabelValues(topic, "failed", err.Error()).Inc()
	}

	p.AsyncProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(v),
	}
}
