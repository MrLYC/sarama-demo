package config

// Kafka : Kafka configuration
type Kafka struct {
	Level string `yaml:"level"`
}

// Init : init Kafka
func (l *Kafka) Init() {
	l.Level = "info"
}
