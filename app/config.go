package app

import (
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Env                string        `config:"env"`
	LogJson            bool          `config:"log.json"`
	LogLevel           string        `config:"log.level"`
	ServerPort         string        `config:"server.port"`
	Secret             string        `config:"secret" json:"-"`
	TerminationTimeout time.Duration `config:"termination.timeout"`
}

func DefaultConfig() Config {
	return Config{
		ServerPort:         "8080",
		LogLevel:           "info",
		TerminationTimeout: 30 * time.Second,
	}
}

func LoadConfig() Config {
	const msg = "Failed to load application config"

	var tag = "config"
	var k = koanf.New(".")

	err := k.Load(structs.Provider(DefaultConfig(), tag), nil)
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}

	err = k.Load(env.Provider("MIND_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "MIND_")), "_", ".", -1)
	}), nil)
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}

	config := Config{}
	err = k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: tag})
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}

	return config
}
