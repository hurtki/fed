package ollama

type GenerateRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Format  string                 `json:"format"`
	Options GenerateRequestOptions `json:"options"`
}

type GenerateRequestOptions struct {
	Temperature float32 `json:"temperature"`
}
type GenerateResponse struct {
	Response string `json:"response"`
}
