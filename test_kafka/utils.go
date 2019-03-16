package test_kafka

import (
	"test/logging"
)

func checkError(err error) {
	if err == nil {
		return
	}
	logger := logging.GetLogger()
	logger.Fatalf("error %v", err)
}
