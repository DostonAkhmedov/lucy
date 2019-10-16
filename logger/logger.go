package logger

import (
	"github.com/iktakahiro/slclogger"
	"log"
	"os"
)

var defaultLogger *log.Logger

func ToFile(file string, prefix string) *log.Logger {

	if defaultLogger != nil {
		return defaultLogger
	}

	folderPath := "logs/"
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	filepath := folderPath + file
	log.Println(filepath)
	f, err := os.OpenFile(filepath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	defaultLogger := log.New(f, prefix, log.LstdFlags)

	return defaultLogger
}

func SlcLogger(webHookUrl string) *slclogger.SlcLogger {
	slacklogger, _ := slclogger.NewSlcLogger(&slclogger.LoggerParams{
		WebHookURL:     webHookUrl,
		LogLevel:       slclogger.LevelDebug,
		DefaultChannel: "",
		DebugChannel:   "",
		InfoChannel:    "",
		WarnChannel:    "",
		ErrorChannel:   "",
		UserName:       "",
		DefaultTitle:   "",
		IconURL:        "",
	})

	return slacklogger
}
