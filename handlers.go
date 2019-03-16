package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"test_kafka/config"
	"test_kafka/logging"
)

func initRandomSeed() bool {
	rand.Seed(time.Now().Unix())
	return true
}

func initConfiguration() bool {
	var err error
	err = config.Configuration.Read()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	err = config.Configuration.Validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return false
	}
	return true
}

func initLogging() bool {
	logging.Init()
	return true
}
