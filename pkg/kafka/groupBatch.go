package kafka

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"thumb/model"
	"time"

	"github.com/IBM/sarama"
)

type ConsumeHandler interface {
	Setup(sarama.ConsumerGroupSession) error
	Cleanup(sarama.ConsumerGroupSession) error
	ConsumeClaim(sarama.ConsumerGroupSession, sarama.ConsumerGroupClaim) error
}

type LikeHandler struct {
	ready       chan bool
	message     chan *sarama.ConsumerMessage
	batchSize   int
	flushTime   time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	mu          sync.Mutex
	batchBuffer []*sarama.ConsumerMessage
	lastCommit  time.Time
}

var _ ConsumeHandler = (*LikeHandler)(nil)

func (l *LikeHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Printf("Consumer group session setup for member: %s\n", session.MemberID())
	l.wg.Add(1)
	go l.processBatches(session)
	close(l.ready)
	return nil
}

func (l *LikeHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Consumer group session cleanup for member: ", session.MemberID())
	l.cancel()
	l.wg.Wait()
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.batchBuffer) > 0 {
		log.Println("Processing remaining messages in buffer during cleanup...")
		l.commit(session, l.batchBuffer)
		l.batchBuffer = nil
	}
	return nil
}

func (l *LikeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Println("Consumer group session claim for member: ", session.MemberID())
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				log.Println("Claim messages channel closed")
				return nil
			}
			select {
			case l.message <- msg:
			case <-session.Context().Done():
				log.Printf("Session context cancelled for topic: %s, partition: %d, during message send.\n", claim.Topic(), claim.Partition())
				return nil
			}
		case <-session.Context().Done():
			log.Printf("Session context cancelled for topic: %s, partition: %d, during message send.\n", claim.Topic(), claim.Partition())
			return nil
		}
	}
}

func NewLikeHandler(batchSize int, flushTime time.Duration) *LikeHandler {
	ctx, cancel := context.WithCancel(context.Background())
	return &LikeHandler{
		ready:      make(chan bool),
		message:    make(chan *sarama.ConsumerMessage, 10000),
		batchSize:  batchSize,
		flushTime:  flushTime,
		ctx:        ctx,
		cancel:     cancel,
		lastCommit: time.Now(),
	}
}

func (l *LikeHandler) processBatches(session sarama.ConsumerGroupSession) {
	defer l.wg.Done()
	ticker := time.NewTicker(l.flushTime)
	defer ticker.Stop()

	for {
		select {
		case msg := <-l.message:

			var curBatch []*sarama.ConsumerMessage
			l.mu.Lock()
			l.batchBuffer = append(l.batchBuffer, msg)

			if len(l.batchBuffer) >= l.batchSize {
				log.Printf("Batch size reached (%d), processing batch.\n", len(l.batchBuffer))
				curBatch = l.batchBuffer
				l.batchBuffer = nil
				l.lastCommit = time.Now()
			}
			l.mu.Unlock()

			if curBatch != nil {
				l.commit(session, curBatch)
			}
		case <-ticker.C:
			var curBatch []*sarama.ConsumerMessage
			l.mu.Lock()
			if len(l.batchBuffer) > 0 && time.Since(l.lastCommit) > l.flushTime {
				log.Printf("Batch interval reached, processing batch (%d messages).\n", len(l.batchBuffer))
				curBatch = l.batchBuffer
				l.batchBuffer = nil
				l.lastCommit = time.Now()
			}
			l.mu.Unlock()
			if curBatch != nil {
				l.commit(session, curBatch)
			}
		case <-l.ctx.Done():
			log.Println("Batch processing goroutine shutting down.")
			return
		case <-session.Context().Done():
			log.Println("Session context cancelled, batch processing goroutine shutting down.")
			return
		}
	}
}

func (l *LikeHandler) commit(session sarama.ConsumerGroupSession, batch []*sarama.ConsumerMessage) {
	counts := make(map[string]int)
	for _, topic := range batch {
		var model model.SocialAction
    if err := json.Unmarshal(topic.Value, &model); err != nil {
  		log.Println("error in json unmarshal")
		}
		counts[model.PostID] += model.Like
	}

	//	TODO:write in redis
	log.Printf("Simulating writing %d like counts to database: %v\n", len(counts), counts["1"])

	last := batch[len(batch)-1]
	session.MarkMessage(last, "")
	session.Commit()
	log.Printf("Committed offsets for batch up to Topic=%s, Partition=%d, Offset=%d\n",
		last.Topic, last.Partition, last.Offset)
}
