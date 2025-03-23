package tasks

type Task struct {
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Input       string   `json:"input"`
    Output      string   `json:"output"`
    Tags        []string `json:"tags"`
}
