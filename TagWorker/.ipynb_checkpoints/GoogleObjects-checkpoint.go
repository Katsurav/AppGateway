package tagworker

type PalmRequest struct {
	Question string `json:"question"`
	Model    string `json:"model"`
	Pre      string `json:"pre_prompt"`
	Post     string `json:"post_prompt"`
}

type PalmResponse struct {
	Answer    string          `json:"response"`
	Model     string          `json:"model"`
	Citations []PalmCitations `json:"Document"`
}

type PalmCitations struct {
	Title string `json:"title"`
}

type ESRequest struct {
	Question  string `json:"query"`
	Page_Size int    `json:"page_size"`
	Offset    int    `json:"offset"`
}

type ESResponse struct {
	Results []struct {
		Id       string `json:"id"`
		Document struct {
			Name              string `json:"name"`
			Id                string `json:"id"`
			DerivedStructData struct {
				ExtractiveAnswers []struct {
					Content    string `json:"content"`
					PageNumber string `json:"pageNumber,omitempty"`
				} `json:"extractive_answers"`
				Link string `json:"link"`
			} `json:"derivedStructData"`
		} `json:"document"`
	} `json:"results"`
	TotalSize          int    `json:"totalSize"`
	AttributionToken   string `json:"attributionToken"`
	NextPageToken      string `json:"nextPageToken"`
	GuidedSearchResult struct {
	} `json:"guidedSearchResult"`
	Summary struct {
	} `json:"summary"`
}
