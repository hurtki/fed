package config

import (
	"github.com/DeanPDX/dotconfig"
)

type GeminiConfig struct {
	Token string `env:"GEMINI_TOKEN,required"`
	Model string `env:"GEMINI_MODEL,requried"`
}

func LoadGeminiConfig() (GeminiConfig, error) {
	cfg, err := dotconfig.FromFileName[GeminiConfig](".env")
	if err != nil {
		return GeminiConfig{}, err
	}
	return cfg, nil
}
