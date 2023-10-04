package tagworker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type AzureLLM struct {
	Id               string
	EngineName       string
	EngineUrl        string
	EnginePrePrompt  string
	EnginePostPrompt string
	EngineTimeout    int64

	// Password
	EngineToken string

	IsCompleted bool
	DidError    string

	// Callback (Engine string, Response string, Citations []string, DidError bool, Runtime int64) ()
	Callback WorkerCallback

	Query         string
	Response      string
	QueryResponse string
	Citations     []string
	Runtime       int64
	Model         string

	Request AzureRequest
}

func (l *AzureLLM) Init(id string, model string) {
	l.Id = id
	l.EngineName = "Azure"
	l.EngineUrl = "https://knowitall-webapp.azurewebsites.net/api_vik"
	l.EnginePrePrompt = "Image you are a Support Bot, your goal is to provide accurate and helpful detail information about the question based on the information in the document"
	l.EnginePostPrompt = "Please summarize your findings, and answer I don't know if you cannot find anything"
	l.EngineTimeout = 3500
	l.EngineToken = ""
	l.Callback = nil
	//l.Model = "gpt-35-turbo"
	l.Model = model
}

func (l *AzureLLM) LLMEngineName() string {
	return l.EngineName
}

func (l *AzureLLM) LLMEngineUrl() string {
	return l.EngineUrl
}

func (l *AzureLLM) LLMEnginePrePrompt() string {
	return l.EnginePrePrompt
}

func (l *AzureLLM) LLMEnginePostPrompt() string {
	return l.EnginePostPrompt
}

func (l *AzureLLM) LLMEngineModel() string {
	return l.Model
}

func (l *AzureLLM) Completed() bool {
	return l.IsCompleted
}

func (l *AzureLLM) Error() string {
	return l.DidError
}

func (l *AzureLLM) LLMCitations() []string {
	return l.Citations
}

func (l *AzureLLM) LLMResponse() string {
	if len(l.DidError) > 0 {
		return l.DidError
	} else {
		return l.QueryResponse
	}
}

func (l *AzureLLM) LLMRuntime() int64 {
	return l.Runtime
}

func (l *AzureLLM) Assign(query string, timeout int64, pre string, post string, callback WorkerCallback) {
	l.Callback = callback
	l.Query = query

	if timeout > 0 {
		l.EngineTimeout = timeout
	}

	if len(pre) > 0 {
		if strings.HasPrefix(pre, "!") {
			l.EnginePrePrompt = ""
			pre = pre[1:]
		}
		l.EnginePrePrompt = pre + " " + l.EnginePrePrompt
	}

	if len(post) > 0 {
		if strings.HasPrefix(post, "!") {
			l.EnginePostPrompt = ""
			post = post[1:]
		}
		l.EnginePostPrompt = l.EnginePostPrompt + " " + post
	}
}

func (l *AzureLLM) DoCallback() {
	isError := len(l.DidError) > 0
	if isError {
		l.Callback(l.Id, l.EngineName, l.DidError, l.Citations, isError, l.Runtime)
	} else {
		l.Callback(l.Id, l.EngineName, l.Response, l.Citations, isError, l.Runtime)
	}
}

func (l *AzureLLM) Run(wg *sync.WaitGroup) {

	now := time.Now()
	defer func() {
		l.Runtime = time.Since(now).Milliseconds()
		l.DoCallback()
		wg.Done()
	}()

	l.Request = AzureRequest{}
	l.Request.Pre = l.EnginePrePrompt
	l.Request.Post = l.EnginePostPrompt
	l.Request.Question = l.Query
	l.Request.Model = l.Model

	marshalled, err := json.Marshal(l.Request)
	if err != nil {
		l.DidError = fmt.Sprintf("impossible to marshall request: %s", err)
		log.Printf(l.DidError, err)
		return
	}

	log.Printf("Request Body: %s", string(marshalled))

	req, err := http.NewRequest("POST", l.EngineUrl, bytes.NewReader(marshalled))
	if err != nil {
		l.DidError = fmt.Sprintf("impossible to build request: %s", err)
		log.Printf(l.DidError, err)
		return
	}

	// add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-functions-key", l.EngineToken)

	client := http.Client{Timeout: 60 * time.Second}
	// send the request
	res, err := client.Do(req)
	if err != nil {
		l.DidError = fmt.Sprintf("impossible to send request: %s", err)
		log.Printf(l.DidError, err)
		return
	}
	log.Printf("status Code: %d", res.StatusCode)

	// we do not forget to close the body to free resources
	// defer will execute that at the end of the current function
	defer res.Body.Close()
	// read body
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		l.DidError = fmt.Sprintf("impossible to read all body of response: %s", err)
		log.Printf(l.DidError, err)
		return
	}
	log.Printf("res body: %s", string(resBody))

	l.Response = string(resBody)

	result := AzureResponse{}

	// now we need to decode the string back to a response we can use
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		l.DidError = fmt.Sprintf("impossible to decode  response: %s", err)
		log.Printf(l.DidError, err)
		return
	}

	l.QueryResponse = result.Answer
	l.Citations = []string{}
	for _, item := range result.Citations {
		l.Citations = append(l.Citations, item.Title)
	}

}
