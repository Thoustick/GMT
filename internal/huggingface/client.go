package huggingface

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Thoustick/GMT/pkg/logger"
	"github.com/Thoustick/GMT/internal/config"
)

// HuggingFaceClient - клиент для Hugging Face API
type HuggingFaceClient struct {
	apiKey string
	model  string
	client *http.Client
	log    logger.Logger
}

// NewClient - создание клиента Hugging Face
func NewClient(cfg *config.Config, log logger.Logger) (*HuggingFaceClient, error) {
	if cfg.HugFaceApiKey == "" || cfg.HugFaceModel == "" {
		return nil, errors.New("отсутствуют параметры Hugging Face в конфигурации")
	}

	return &HuggingFaceClient{
		apiKey: cfg.HugFaceApiKey,
		model:  cfg.HugFaceModel,
		client: &http.Client{Timeout: 15 * time.Second},
		log:    log,
	}, nil
}

func (c *HuggingFaceClient) SendRequest(ctx context.Context, prompt string) ([]byte, error) {
    url := "https://api-inference.huggingface.co/models/" + c.model
    requestBody, err := json.Marshal(map[string]string{"inputs": prompt})
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("ошибка Hugging Face API: " + string(body))
    }

    return body, nil // Возвращаем "сырой" JSON-ответ
}
