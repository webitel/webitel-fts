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
	Insecure  bool
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
	Import           ImportConfig
}

type SetupConfig struct {
	Template string
}

type ImportConfig struct {
	Index       string
	Query       string
	File        string
	ColDomainId string
	ColId       string
}
