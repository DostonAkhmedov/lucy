package config

type Config struct {
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
}

func (c Config) GetDBHost() string {
	return c.dbHost
}

func (c Config) GetDBPort() int {
	return c.dbPort
}

func (c Config) GetDBUser() string {
	return c.dbUser
}

func (c Config) GetDBPassword() string {
	return c.dbPassword
}

func (c Config) GetDBName() string {
	return c.dbName
}
