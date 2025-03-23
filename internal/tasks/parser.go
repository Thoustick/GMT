package tasks

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"fmt"
)

// APIResponse - структура ответа Hugging Face
type APIResponse struct {
	GeneratedText string `json:"generated_text"`
}

func (t *Task) ToMarkdownV2() string {
    // Helper function to escape MarkdownV2 special characters
    escape := func(s string) string {
        special := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}
        for _, ch := range special {
            s = strings.ReplaceAll(s, ch, "\\"+ch)
        }
        return s
    }
    
    md := fmt.Sprintf(
        "*%s*\n\n*Описание:* %s\n\n*Входные данные:* `%s`\n\n*Выходные данные:* `%s`\n\n*Теги:* `%s`",
        escape(t.Title),
        escape(t.Description),
        escape(t.Input),
        escape(t.Output),
        strings.Join(t.Tags, "`, `"),
    )
    return md
}

// func (t *Task) ToHTML() string {
//     // Helper function to escape HTML special characters
//     escape := func(s string) string {
//         s = strings.ReplaceAll(s, "&", "&amp;")
//         s = strings.ReplaceAll(s, "<", "&lt;")
//         s = strings.ReplaceAll(s, ">", "&gt;")
//         return s
//     }
    
//     html := fmt.Sprintf(
//         "<b>%s</b>\n\n<b>Описание:</b> %s\n\n<b>Входные данные:</b> <code>%s</code>\n\n<b>Выходные данные:</b> <code>%s</code>\n\n<b>Теги:</b> <code>%s</code>",
//         escape(t.Title),
//         escape(t.Description),
//         escape(t.Input),
//         escape(t.Output),
//         strings.Join(t.Tags, "</code>, <code>"),
//     )
//     return html
// }

// ParseTask extracts the first valid Task JSON from raw text
func ParseTask(rawText string) (*Task, error) {
    // Use a regex that finds complete JSON objects (supports nesting)
    re := regexp.MustCompile(`(?s)\{(?:[^{}]|(?:\{[^{}]*\}))*\}`)
    matches := re.FindAllString(rawText, -1)
    
    if len(matches) == 0 {
        return nil, errors.New("не найден JSON в ответе")
    }
    
    // Try to parse each match until we find a valid task
    for _, match := range matches {
        var task Task
        if err := json.Unmarshal([]byte(match), &task); err == nil {
            // Verify this is a valid task with necessary fields
            if task.Title != "" && task.Description != "" {
                return &task, nil
            }
        }
    }
    
    return nil, errors.New("не найден валидный JSON задачи в ответе")
}

// parseHuggingFaceResponse extracts the generated text from API response
func parseHuggingFaceResponse(rawResponse []byte) (string, error) {
    // First try as a single response
    var singleResponse APIResponse
    if err := json.Unmarshal(rawResponse, &singleResponse); err == nil && singleResponse.GeneratedText != "" {
        return singleResponse.GeneratedText, nil
    }
    
    // Then try as an array of responses
    var multipleResponses []APIResponse
    if err := json.Unmarshal(rawResponse, &multipleResponses); err != nil {
        return "", fmt.Errorf("failed to parse API response: %w", err)
    }
    
    if len(multipleResponses) == 0 || multipleResponses[0].GeneratedText == "" {
        return "", errors.New("пустой ответ от Hugging Face API")
    }
    
    return multipleResponses[0].GeneratedText, nil
}

// cleanResponse removes the original prompt from the generated text.
func cleanResponse(response, prompt string) string {
	cleaned := strings.TrimPrefix(response, prompt)
	return strings.TrimSpace(cleaned)
}