package kafka

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kzerolog"

	"course-service/pkg/logger"
)

type Handler func(ctx context.Context, records []*kgo.Record) error
type ErrorHandler func(ctx context.Context, records []*kgo.Record, err error)

type Consumer struct {
	client       *kgo.Client
	topic        string
	handler      Handler
	errorHandler ErrorHandler
	logger       *zerolog.Logger
}

type ConsumerConfig struct {
	Brokers      []string
	GroupID      string
	Topic        string
	Handler      Handler
	ErrorHandler ErrorHandler
	Logger       *zerolog.Logger
}

func (c ConsumerConfig) validate() error {
	if len(c.Brokers) == 0 {
		return fmt.Errorf("at least one broker is required")
	}
	if c.GroupID == "" {
		return fmt.Errorf("group id is required")
	}
	if c.Topic == "" {
		return fmt.Errorf("topic is required")
	}
	if c.Handler == nil {
		return fmt.Errorf("handler is required")
	}
	return nil
}

func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	kafkaLogger := cfg.Logger.Level(zerolog.InfoLevel).With().Str("module", "consumer").Logger()

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ConsumerGroup(cfg.GroupID),
		kgo.ConsumeTopics(cfg.Topic),
		kgo.DisableAutoCommit(),
		kgo.MetadataMaxAge(60 * time.Second),
		kgo.FetchMaxBytes(50_000_000),
		kgo.FetchMaxPartitionBytes(50_000_000),
		kgo.FetchMaxWait(1 * time.Second),
		kgo.WithLogger(kzerolog.New(new(kafkaLogger))),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	return &Consumer{
		client:       client,
		topic:        cfg.Topic,
		handler:      cfg.Handler,
		errorHandler: cfg.ErrorHandler,
		logger:       cfg.Logger,
	}, nil
}

func (c *Consumer) Topic() string {
	return c.topic
}

func (c *Consumer) Run(ctx context.Context) {
	c.logger.
		Info().
		Str("topic", c.topic).
		Msg("starting consumer")

	for {
		fetches := c.client.PollFetches(ctx)

		if fetches.IsClientClosed() || ctx.Err() != nil {
			c.logger.Debug().Msg("stop consuming: client was closed or context was cancelled")
			return
		}

		if errs := fetches.Errors(); len(errs) > 0 {
			c.logger.Error().Any("fetches errors", errs).Send()
			continue
		}

		var wg sync.WaitGroup

		fetches.EachPartition(
			func(partition kgo.FetchTopicPartition) {
				records := partition.Records

				if logger.IsTrace(c.logger) {
					c.logger.
						Trace().
						Msgf(
							"topic: %s, partition: %d, fetched %d records",
							partition.Topic,
							partition.Partition,
							len(records),
						)
				}

				if len(records) == 0 {
					return
				}

				wg.Add(1)
				go c.handleRecords(ctx, &wg, records)
			},
		)

		wg.Wait()
	}
}

func (c *Consumer) handleRecords(ctx context.Context, wg *sync.WaitGroup, records []*kgo.Record) {
	defer wg.Done()

	defer c.recoverHandler()

	if err := c.handler(ctx, records); err != nil {
		if c.errorHandler != nil {
			c.errorHandler(ctx, records, err)
		}
		return
	}

	commitCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.client.CommitRecords(commitCtx, records...); err != nil {
		c.logger.
			Error().
			Err(err).
			Str("topic", c.topic).
			Int32("partition", records[0].Partition).
			Int64("offset", records[0].Offset).
			Msg("commit records")
	}
}

func (c *Consumer) recoverHandler() {
	if r := recover(); r != nil {
		c.logger.Error().
			Any("panic", r).
			Bytes("stacktrace", debug.Stack()).
			Msg("recovered from panic")
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}
