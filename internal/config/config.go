package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string   `yaml:"env" env:"ENV" env-required:"true"`
	Token      string   `yaml:"token" env:"TOKEN" env-required:"true" yaml-required:"true"`
	ListenAddr string   `yaml:"listen_addr" env:"LISTEN_ADDR" env-default:"0.0.0.0:8080"`
	Postgres   Postgres `yaml:"database"`
	Webhook    Webhook  `yaml:"webhook"`
}

type Postgres struct {
	Dsn string `yaml:"dsn" env:"POSTGRES_DSN" env-required:"true"`
}

type Webhook struct {
	Secret string `yaml:"secret" env:"WEBHOOK_SECRET" env-required:"true"`
	Domain string `yaml:"domain" env:"WEBHOOK_DOMAIN" env-required:"true"`
}

type ConfigOptions struct {
	ConfigPath string
}

// Functions that start with the Must prefix require that the config is loaded, otherwise panic will be thrown.
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
		// log.Fatal("CONFIG_PATH is not set")

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
