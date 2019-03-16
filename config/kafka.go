package config

import (
	"time"
)

// Kafka : Kafka configuration
type Kafka struct {
	Brokers        string        `yaml:"brokers"`
	SessionTimeout time.Duration `yaml:"session_timeout"`
}

// Init : init Kafka
func (k *Kafka) Init() {
	k.Brokers = ":9092"
	k.SessionTimeout = 10 * time.Second
}
