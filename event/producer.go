package event

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	BatchProduceFeedEvent(ctx context.Context, event []FeedEvent) error
}

type SaramaProducer struct {
	producer sarama.SyncProducer
}

func NewSaramaProducer(producer sarama.SyncProducer) Producer {
	return &SaramaProducer{producer: producer}
}

func (p *SaramaProducer) BatchProduceFeedEvent(ctx context.Context, event []FeedEvent) error {
	msgs := make([]*sarama.ProducerMessage, 0, len(event))
	for _, e := range event {
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}
		msgs = append(msgs, &sarama.ProducerMessage{
			Topic: topicFeedEvent,
			Value: sarama.ByteEncoder(data),
		})
	}
	return p.producer.SendMessages(msgs)
}
