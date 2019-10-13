package logger

import (
	"fmt"
	"github.com/DostonAkhmedov/lucy/config"
	"github.com/iktakahiro/slclogger"
	"log"
	"os"
)

var defaultLogger *log.Logger

func ToFile(file string, prefix string) *log.Logger {

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

func SlcLogger(conf *config.Config) *slclogger.SlcLogger {
	if conf == nil {
		conf = config.Init()
	}
	slacklogger, _ := slclogger.NewSlcLogger(&slclogger.LoggerParams{
		WebHookURL:     conf.GetSlcWebHook(),
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
