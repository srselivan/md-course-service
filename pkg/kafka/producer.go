package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kzerolog"
)

type Producer struct {
	client *kgo.Client
	logger *zerolog.Logger
}

type ProducerConfig struct {
	Brokers []string
	Logger  *zerolog.Logger
}

func NewProducer(cfg ProducerConfig) (*Producer, error) {
	kafkaLogger := cfg.Logger.Level(zerolog.InfoLevel).With().Str("module", "producer").Logger()

	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Brokers...),
		kgo.ProducerBatchMaxBytes(16_000_000),
		kgo.ProduceRequestTimeout(10 * time.Second),
		kgo.MetadataMaxAge(60 * time.Second),
		kgo.MaxBufferedRecords(1_000_000),
		kgo.WithLogger(kzerolog.New(new(kafkaLogger))),
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("new client: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = client.Ping(pingCtx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return &Producer{
		client: client,
		logger: cfg.Logger,
	}, nil
}

type ProduceOpts struct {
	Key     []byte
	Value   []byte
	Topic   string
	Headers []kgo.RecordHeader
}

func (p *Producer) ProduceSync(ctx context.Context, opts ProduceOpts) error {
	record := &kgo.Record{
		Topic:   opts.Topic,
		Key:     opts.Key,
		Value:   opts.Value,
		Headers: opts.Headers,
	}
	results := p.client.ProduceSync(ctx, record)
	if err := results.FirstErr(); err != nil {
		return fmt.Errorf("produce: %w", err)
	}
	return nil
}

func (p *Producer) Produce(ctx context.Context, opts ProduceOpts) error {
	record := &kgo.Record{
		Topic:   opts.Topic,
		Key:     opts.Key,
		Value:   opts.Value,
		Headers: opts.Headers,
	}
	p.client.Produce(ctx, record, p.producePromise)
	return nil
}

func (p *Producer) producePromise(record *kgo.Record, err error) {
	if err == nil {
		return
	}
	p.logger.
		Error().
		Err(err).
		Str("topic", record.Topic).
		Str("key", string(record.Key)).
		Str("value", string(record.Value)).
		Send()
}

func (p *Producer) Ping(ctx context.Context) error {
	return p.client.Ping(ctx)
}

func (p *Producer) Close() {
	p.client.Close()
}
