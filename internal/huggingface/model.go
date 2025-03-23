package huggingface

type GPTRequest struct {
    Model    string `json:"model"`
    Messages []GPTMessage `json:"messages"`
    Temperature float64 `json:"temperature"`
}

type GPTMessage struct {
        Role    string `json:"role"`
        Content string `json:"content"`
}

type GPTResponse struct {
    Choices []struct {
        Message struct {
            Content string `json:"content"`
        } `json:"message"`
    } `json:"choices"`
}