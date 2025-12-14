package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env:"LOGBOT_ENV" env-required:"true"`
	Token      string     `yaml:"token" env:"LOGBOT_TOKEN" env-required:"true" yaml-required:"true"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"database"`
	Webhook    Webhook    `yaml:"webhook"`
	Superuser  string     `yaml:"superuser" env:"LOGBOT_SUPERUSER" env-required:"true" yaml-required:"true"`
	LogCleaner LogCleaner `yaml:"log_cleaner"`
}

type LogCleaner struct {
	Interval time.Duration `yaml:"interval" env:"LOGBOT_LOG_CLEANER_INTERVAL" env-default:"10m"`
	// 672 hours is 4 weeks
	Retain time.Duration `yaml:"retain" env:"LOGBOT_LOG_CLEANER_RETAIN" env-default:"672h"`
}

type Postgres struct {
	Dsn string `yaml:"dsn" env:"POSTGRES_DSN" env-required:"true"`
}

type HTTPServer struct {
	Host         string        `yaml:"host" env:"LOGBOT_HOST" env-default:"0.0.0.0"`
	Port         string        `yaml:"port" env:"LOGBOT_PORT" env-default:"80"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env:"LOGBOT_READ_TIMEOUT" env-required:"true"`
	WriteTimeout time.Duration `yaml:"write_timeout" env:"LOGBOT_WRITE_TIMEOUT" env-required:"true"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env:"LOGBOT_IDLE_TIMEOUT" env-required:"true"`
}

type Webhook struct {
	Secret string `yaml:"secret" env:"LOGBOT_WEBHOOK_SECRET" env-required:"true"`
	Domain string `yaml:"domain" env:"LOGBOT_WEBHOOK_DOMAIN" env-required:"true"`
}

type ConfigOptions struct {
	ConfigPath string
}

func MustLoad(opts *ConfigOptions) *Config {
	var (
		cfg        Config
		configPath string
	)

	if opts != nil {
		configPath = opts.ConfigPath
	}

	if configPath == "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	if configPath != "" {
		// check if file exists.
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			log.Fatalf("config file %s does not exist", configPath)
		}

		// Read from file.
		if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
			log.Fatalf("failed to load config from %s: %v", configPath, err)
		}
	}

	// Load configuration from the environment, overriding any previously loaded config file values.
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("failed to load config (neither yaml nor env are defined): %v", err)
	}

	return &cfg
}
