package config

import "github.com/namsral/flag"

var defaultConfig *Config

func Init() *Config {

	if defaultConfig != nil {
		return defaultConfig
	}

	defaultConfig := &Config{}

	flag.StringVar(&defaultConfig.dbHost, "db-host", "localhost", "Database host")
	flag.IntVar(&defaultConfig.dbPort, "db-port", 3306, "Database port")
	flag.StringVar(&defaultConfig.dbUser, "db-user", "root", "Database username")
	flag.StringVar(&defaultConfig.dbPassword, "db-password", "", "Database password")
	flag.StringVar(&defaultConfig.dbName, "db-name", "lucy", "Database name")

	flag.StringVar(&defaultConfig.slcWebhook, "slc-webhook", "https://hooks.slack.com/services/YOUR_WEBHOOK_URL", "Slack Logger")

	flag.Parse()

	return defaultConfig
}
