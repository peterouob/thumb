package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"math/rand"
	"sync"
	"thumb/configs"
	"time"
)

type LikeConsumer struct {
	consumer   sarama.ConsumerGroup
	groupId    string
	topic      []string
	batchSize  int
	flushTime  time.Duration
	handler    *LikeHandler
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	backoffCfg *configs.BackOffConfig
	random     *rand.Rand
	sync.Mutex
}

func NewConsumer(groupId string, topic []string, batchSize int, flushTime time.Duration) *LikeConsumer {
	conf := sarama.NewConfig()
	conf.Consumer.Fetch.Default = 1024 * 1024
	conf.Consumer.MaxWaitTime = 500 * time.Millisecond
	conf.Consumer.Retry.Backoff = 500 * time.Millisecond
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Consumer.Offsets.AutoCommit.Enable = false

	consumer, err := sarama.NewConsumerGroup([]string{"192.168.0.100:9092"}, groupId, conf)
	if err != nil {
		log.Fatalln(errors.New("error in NewConsumer :" + err.Error()))
	}

	ctx, cancel := context.WithCancel(context.Background())

	c := &LikeConsumer{
		consumer:   consumer,
		groupId:    groupId,
		topic:      topic,
		batchSize:  batchSize,
		flushTime:  flushTime,
		ctx:        ctx,
		cancel:     cancel,
		handler:    NewLikeHandler(batchSize, flushTime),
		backoffCfg: configs.NewBackOffConfig(true),
		random:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	go c.handleError()

	log.Println("init kafka consumer ...")
	return c
}

func (c *LikeConsumer) StartConsume() {
	c.wg.Add(1)
	go func() {
		for {
			if err := c.consumer.Consume(c.ctx, c.topic, c.handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					log.Println("Consumer group closed, exiting consume loop.")
					return
				}

				log.Printf("Error from consumer group: %v, sarting backoff ...\n", err)
				time.Sleep(c.backOff())
			}

			if c.ctx.Err() != nil {
				log.Println("Context cancelled, exiting consume loop.")
				return
			}

			c.handler.ready = make(chan bool)
		}
	}()

	<-c.handler.ready
	log.Printf("Sarama batch consumer group up and running! [topic: %s, group: %s]\n", c.topic, c.groupId)
}

func (c *LikeConsumer) handleError() {
	for err := range c.consumer.Errors() {
		if err != nil {
			// TODO:prometheus
			log.Println(err.Error())
		}
	}
}

func (c *LikeConsumer) Close() error {
	log.Println("Closing Kafka consumer group...")
	c.cancel()
	c.wg.Wait()
	if err := c.consumer.Close(); err != nil {
		return fmt.Errorf("error closing consumer group client: %w", err)
	}
	log.Println("Kafka consumer group closed.")
	return nil
}

func (c *LikeConsumer) backOff() time.Duration {
	curRetries := c.backoffCfg.GetRetries()

	if c.backoffCfg.GetRetries() == 0 {
		return c.backoffCfg.GetBaseDelay()
	}

	curBackOff, maxTime := float64(c.backoffCfg.GetBaseDelay()), float64(c.backoffCfg.GetMaxDelay())

	for curBackOff < maxTime && curRetries > 0 {
		curBackOff = curBackOff * c.backoffCfg.GetMultiplier()
		curRetries--
	}

	if curBackOff > maxTime {
		curBackOff = maxTime
	}

	curBackOff *= 1.0 + c.backoffCfg.GetJitter()*(c.random.Float64()*2-1)
	if curBackOff < 0 {
		return 0
	}
	return time.Duration(curBackOff)
}
