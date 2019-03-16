package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/subcommands"

	"github.com/mrlyc/sarama-demo/config"
	"github.com/mrlyc/sarama-demo/sarama"
)

type initialHandler func() bool

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")
	subcommands.Register(&config.VersionCommand{}, "")
	subcommands.Register(&config.ConfInfoCommand{}, "")
	subcommands.Register(&sarama.Productor{}, "")
	subcommands.Register(&sarama.Consumer{}, "")

	flag.StringVar(
		&(config.Configuration.ConfigurationPath),
		"c", config.Configuration.ConfigurationPath,
		"Configuration file",
	)

	flag.Parse()

	initialHandlers := []initialHandler{
		initRandomSeed,
		initConfiguration,
	}

	for _, handler := range initialHandlers {
		if !handler() {
			os.Exit(255)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			switch <-signalCh {
			case syscall.SIGINT:
				cancel()
			}
		}
	}()

	os.Exit(int(subcommands.Execute(ctx)))
}
