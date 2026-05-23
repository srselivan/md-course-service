package courses

import (
	"context"
	"encoding/json"

	"course-service/pkg/kafka"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	fileTopic = "file-topic"

	EventFileLoaded  = "FileLoadedEvent"
	EventFileDeleted = "FileDeletedEvent"
)

type FilesClient struct {
	producer *kafka.Producer
	logger   *zerolog.Logger
}

func NewFilesClient(producer *kafka.Producer, logger *zerolog.Logger) *FilesClient {
	return &FilesClient{
		producer: producer,
		logger:   logger,
	}
}

func (c *FilesClient) SendEvent(ctx context.Context, fileId, eventType string) error {
	if c.producer == nil {
		return nil
	}

	body, err := json.Marshal(map[string]string{"fileId": fileId})
	if err != nil {
		c.logger.Error().Err(err).
			Str("eventType", eventType).
			Str("fileId", fileId).
			Msg("failed to marshal event")
		return err
	}

	if err = c.producer.Produce(ctx, kafka.ProduceOpts{
		Topic: fileTopic,
		Value: body,
		Headers: []kgo.RecordHeader{
			{Key: "event_type", Value: []byte(eventType)},
		},
	}); err != nil {
		c.logger.Error().Err(err).
			Str("eventType", eventType).
			Str("fileId", fileId).
			Msg("failed to send file event to kafka")
		return err
	}

	return nil
}
