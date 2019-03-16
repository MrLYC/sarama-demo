package sarama

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/subcommands"

	"github.com/mrlyc/sarama-demo/config"
	"github.com/mrlyc/sarama-demo/logging"
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

	client, err := sarama.NewClient(split(config.Configuration.Kafka.Brokers), conf)
	checkError(err)

	partitions, err := client.Partitions(p.topic)
	checkError(err)

	producer, err := sarama.NewSyncProducerFromClient(client)
	checkError(err)

	interval, err := time.ParseDuration(p.interval)
	checkError(err)

	ticker := time.NewTicker(interval)

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case t := <-ticker.C:
			data := fmt.Sprintf("%v", t.Unix())
			partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
				Partition: partitions[rand.Intn(len(partitions))],
				Topic:     p.topic,
				Value:     sarama.StringEncoder(data),
			})
			if err != nil {
				logger.Infof("send to %s error %v", p.topic, err)
			} else {
				logger.Infof("send to %s/%d:%d %v", p.topic, partition, offset, data)
			}
		}
	}

	ticker.Stop()
	checkError(producer.Close())
	checkError(client.Close())

	return subcommands.ExitSuccess
}
