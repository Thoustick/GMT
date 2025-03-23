package tasks

import (
	"context"
	"fmt"

	"github.com/Thoustick/GMT/internal/huggingface"
	"github.com/Thoustick/GMT/pkg/logger"

)

// TaskGenerator defines the interface for generating tasks.
type TaskGenerator interface {
	GenerateTask() (string, error)
}

// TaskGeneratorImpl implements TaskGenerator using a HuggingFaceClient.
type TaskGeneratorImpl struct {
	Client *huggingface.HuggingFaceClient
	Log	logger.Logger
}

// NewTaskGeneratorImpl creates a new TaskGeneratorImpl.
func NewTaskGeneratorImpl(client *huggingface.HuggingFaceClient, log logger.Logger) TaskGeneratorImpl {
	return TaskGeneratorImpl{
		Client: client,
		Log: log,
	}
}

func (tg *TaskGeneratorImpl) GenerateTask() (string, error) {
    prompt := `Сгенерируй задачу для обучения языку Goolang в формате JSON на английском:
{
    "title": "Название",
    "description": "Описание",
    "input": "Пример входа",
    "output": "Пример выхода",
    "tags": ["arrays", "dp"]
}`

    // 1. Отправляем запрос в Hugging Face API
    rawResponse, err := tg.Client.SendRequest(context.Background(), prompt)
    if err != nil {
        tg.Log.Error("Hugging Face API request error", err, nil)
        return "", err
    }
    // 2. Парсим JSON-ответ
    generatedText, err := parseHuggingFaceResponse(rawResponse)
    if err != nil {
		tg.Log.Error("JSON parsing error", err, map[string]interface{}{
			"response": string(rawResponse),
		})
	return "", err
    }

    // 3. Убираем оригинальный промпт (если API возвращает его)
    cleanedResponse := cleanResponse(generatedText, prompt)

    // 4. Извлекаем JSON-задачу
    task, err := ParseTask(cleanedResponse)
    if err != nil {
		tg.Log.Error("JSON extracting error", err, map[string]interface{}{
			"response": cleanedResponse,
		})
	return "", fmt.Errorf("can not extract from answer Hugging Face API: %w", err)
    }
	tg.Log.Info("New task was generated", map[string]interface{}{
		"title": task.Title,
	})
    return task.ToMarkdownV2(), nil
}
