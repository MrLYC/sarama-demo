package test_kafka

import (
	"context"
	"flag"

	"github.com/google/subcommands"

	"test_kafka/logging"
)

// Productor :
type Productor struct{}

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
func (*Productor) SetFlags(f *flag.FlagSet) {

}

// Execute :
func (p *Productor) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	logger := logging.GetLogger()
	logger.Infof("productor running")

	return subcommands.ExitSuccess
}
