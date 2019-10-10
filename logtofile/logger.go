package logtofile

import (
	"fmt"
	"log"
	"os"
	"time"
)

var defaultLogger *log.Logger

func Create(file string, prefix string) *log.Logger {

	if defaultLogger != nil {
		return defaultLogger
	}

	year, month, _ := time.Now().Date()
	filepath := fmt.Sprintf("lucy/logs/%d_%d/%s", year, int(month), file)
	f, err := os.OpenFile(filepath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	defaultLogger := log.New(f, prefix, log.LstdFlags)

	return defaultLogger
}
