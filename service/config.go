package service

import (
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
	"strings"
)

type Config struct {
	Env string `config:"env"`
}

func defaults() Config {
	return Config{
		Env: "lcl",
	}
}

func LoadConfig() Config {
	var tag = "config"
	var k = koanf.New(".")

	err := k.Load(structs.Provider(defaults(), tag), nil)
	if err != nil {
		panic(err)
	}

	err = k.Load(env.Provider("MIND_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "MIND_")), "_", ".", -1)
	}), nil)
	if err != nil {
		panic(err)
	}

	config := Config{}
	err = k.UnmarshalWithConf("", &config, koanf.UnmarshalConf{Tag: tag})
	if err != nil {
		panic(err)
	}

	return config
}
