package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/hurtki/fed/internal/config"
)

type OllamaClient struct {
	cfg config.OllamaConfig
}

func NewOllamaClient(cfg config.OllamaConfig) OllamaClient {
	return OllamaClient{
		cfg: cfg,
	}
}

func (c *OllamaClient) Generate(ctx context.Context, prompt string) (string, error) {
	reader, writer := io.Pipe()

	go func() {
		defer writer.Close()

		err := json.NewEncoder(writer).Encode(GenerateRequest{
			Model:  c.cfg.Model,
			Prompt: prompt,
		})

		if err != nil {
			writer.CloseWithError(err)
		}
	}()

	req, err := http.NewRequestWithContext(ctx, "POST", c.cfg.Endpoint+"/generate", reader)
	if err != nil {
		return "", fmt.Errorf("can't create request: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("can't send http request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		bd, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("not success status code from model: %d, body: %s", res.StatusCode, string(bd))
	}

	dto := GenerateResponse{}
	err = json.NewDecoder(res.Body).Decode(&dto)
	if err != nil {
		return "", fmt.Errorf("wrong json response from model")
	}
	return dto.Response, nil
}
