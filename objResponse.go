package main

type objResponse struct {
	RequestId string      `json:"requestid"`
	Query     string      `json:"query"`
	Results   []objResult `json:"results"`
	Output    []string    `json:"output"`
}

type objResult struct {
	LLMID        string   `json:"llmid"`
	LLMName      string   `json:"llmname"`
	LLMResponse  string   `json:"llmresponse"`
	LLMCitations []string `json:"llmcitations"`
	LLMCompute   int64    `json:"llmcompute"`
	DidError     bool     `json:"diderror"`
}
