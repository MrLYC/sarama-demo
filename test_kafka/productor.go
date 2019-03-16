package test_kafka

import (
	"context"
	"flag"
	"fmt"
	"time"

	"test_kafka/config"

	"github.com/Shopify/sarama"
	"github.com/google/subcommands"

	"test_kafka/logging"
)

// Productor :
type Productor struct {
	topic    string
	interval string
}

// Name :
func (*Productor) Name() string {
	return "productor"
}

// Synopsis :
func (*Productor) Synopsis() string {
	return "productor"
}

// Usage :
func (*Productor) Usage() string {
	return `productor`
}

// SetFlags :
func (p *Productor) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.topic, "topic", "test", "product topic")
	f.StringVar(&p.interval, "interval", "1s", "product interval")
}

// Execute :
func (p *Productor) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logger := logging.GetLogger()
	logger.Infof("productor running")

	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 10
	conf.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(split(config.Configuration.Kafka.Brokers), conf)
	checkError(err)

	interval, err := time.ParseDuration(p.interval)
	checkError(err)

	topics := split(p.topic)
	ticker := time.NewTicker(interval)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case t := <-ticker.C:
			data := fmt.Sprintf("%v", t.Unix())
			for _, topic := range topics {
				partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
					Topic: topic,
					Value: sarama.StringEncoder(data),
				})
				if err != nil {
					logger.Infof("send to %s error %v", topic, err)
				} else {
					logger.Infof("send to %s/%d:%d %v", topic, partition, offset, data)
				}
			}
		}
	}

	checkError(producer.Close())
	ticker.Stop()

	return subcommands.ExitSuccess
}
