package tagworker

type AzureRequest struct {
	Question string `json:"question"`
	Model    string `json:"model"`
	Pre      string `json:"pre_prompt"`
	Post     string `json:"post_prompt"`
}

type AzureResponse struct {
	Answer    string           `json:"response"`
	Model     string           `json:"model"`
	Citations []AzureCitations `json:"Filename"`
}

type AzureCitations struct {
	Title string `json:"title"`
}
