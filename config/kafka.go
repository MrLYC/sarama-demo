package config

// Kafka : Kafka configuration
type Kafka struct {
	Brokers string `yaml:"brokers"`
}

// Init : init Kafka
func (l *Kafka) Init() {
	l.Brokers = ":9092"
}
