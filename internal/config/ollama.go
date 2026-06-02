package config

import (
	"bytes"
	"os"

	"github.com/DeanPDX/dotconfig"
)

type OllamaConfig struct {
	Endpoint string `env:"OLLAMA_ENDPOINT,required"`
	Model    string `env:"OLLAMA_MODEL,requried"`
}

func LoadOllamaConfigFromEnvFile(path string) (OllamaConfig, error) {
	return dotconfig.FromFileName[OllamaConfig](path)
}

func LoadOllamaConfigFromEnvVariables() (OllamaConfig, error) {
	var buf bytes.Buffer

	for _, e := range os.Environ() {
		buf.WriteString(e)
		buf.WriteByte('\n')
	}

	return dotconfig.FromReader[OllamaConfig](&buf)
}
