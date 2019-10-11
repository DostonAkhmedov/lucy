package logtofile

import (
	"fmt"
	"log"
	"os"
)

var defaultLogger *log.Logger

func Create(file string, prefix string) *log.Logger {

	if defaultLogger != nil {
		return defaultLogger
	}

	filepath := fmt.Sprintf("lucy/logs/%s", file)
	f, err := os.OpenFile(filepath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defaultLogger := log.New(f, prefix, log.LstdFlags)

	return defaultLogger
}
