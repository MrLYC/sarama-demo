package sarama

import (
	"context"
	"flag"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
	"github.com/google/subcommands"

	"github.com/mrlyc/sarama-demo/config"
	"github.com/mrlyc/sarama-demo/logging"
)

// Consumer :
type Consumer struct {
	topic string
	group string
}

// Name :
func (*Consumer) Name() string {
	return "consumer"
}

// Synopsis :
func (*Consumer) Synopsis() string {
	return "consumer"
}

// Usage :
func (*Consumer) Usage() string {
	return `consumer`
}

// SetFlags :
func (p *Consumer) SetFlags(f *flag.FlagSet) {
	f.StringVar(&p.topic, "topic", "test", "consumer topic")
	f.StringVar(&p.group, "group", "test", "consumer group")
}

// Execute :
func (p *Consumer) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logger := logging.GetLogger()
	logger.Infof("consumer running")

	conf := cluster.NewConfig()
	conf.Consumer.Return.Errors = true
	conf.Group.Return.Notifications = true
	conf.ClientID = "consumer"
	conf.Consumer.Offsets.Initial = sarama.OffsetNewest

	client, err := cluster.NewClient(split(config.Configuration.Kafka.Brokers), conf)
	checkError(err)

	consumer, err := cluster.NewConsumerFromClient(client, p.group, split(p.topic))
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case ntf, ok := <-consumer.Notifications():
			if !ok {
				logger.Infof("notification channel closed")
				break loop
			}
			logger.Infof("notification %v, current %v, claimed %v, released %v", ntf.Type, ntf.Current, ntf.Claimed, ntf.Released)
		case e, ok := <-consumer.Errors():
			if !ok {
				logger.Infof("error channel closed")
				break loop
			}
			logger.Errorf("error %v", e)
		case msg, ok := <-consumer.Messages():
			if !ok {
				logger.Infof("message channel closed")
				break loop
			}
			logger.Infof("receive from %s/%d:%d %v", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		}
	}

	return subcommands.ExitSuccess
}
