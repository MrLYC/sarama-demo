package test_kafka

import (
	"strings"

	"test/logging"
)

func checkError(err error) {
	if err == nil {
		return
	}
	logger := logging.GetLogger()
	logger.Fatalf("error %v", err)
}

func split(s string)[]string {
	return strings.Split(s, ",")
}