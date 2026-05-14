package config

import (
	"strings"

	"github.com/DeanPDX/dotconfig"
)

type OllamaConfig struct {
	Endpoint string `env:"OLLAMA_ENDPOINT,required"`
	Model    string `env:"OLLAMA_MODEL,requried"`
}

func LoadOllamaConfig() (OllamaConfig, error) {
	cfg, err := dotconfig.FromFileName[OllamaConfig](".env")
	if err != nil {
		return OllamaConfig{}, err
	}
	cfg.Endpoint = strings.TrimSuffix(cfg.Endpoint, "/")
	return cfg, nil
}
