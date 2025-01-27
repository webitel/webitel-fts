package config

type Service struct {
	Id      string
	Address string
	Consul  string
}

type Pubsub struct {
	Address string
}

type OpenSearchConfig struct {
	Addresses string
	Username  string
	Password  string
}

type LogSettings struct {
	Lvl     string
	Json    bool
	Otel    bool
	File    string
	Console bool
}

type SqlSettings struct {
	DSN string
}

type Config struct {
	Service          Service
	Log              LogSettings
	OpenSearchConfig OpenSearchConfig
	Pubsub           Pubsub
	SqlSettings      SqlSettings
}

type SetupConfig struct {
	Template string
}
