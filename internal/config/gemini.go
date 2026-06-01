package config

import (
	"bytes"
	"os"

	"github.com/DeanPDX/dotconfig"
)

type GeminiConfig struct {
	Token string `env:"GEMINI_TOKEN,required"`
	Model string `env:"GEMINI_MODEL,requried"`
}

func LoadGeminiConfigFromEnvFile(path string) (GeminiConfig, error) {
	return dotconfig.FromFileName[GeminiConfig](path)
}

func LoadGeminiConfigFromEnvVariables() (GeminiConfig, error) {
	var buf bytes.Buffer

	for _, e := range os.Environ() {
		buf.WriteString(e)
		buf.WriteByte('\n')
	}

	return dotconfig.FromReader[GeminiConfig](&buf)
}
