package tagworker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	auth "golang.org/x/oauth2/google"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type GoogleES struct {
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

	Request ESRequest
}

func (l *GoogleES) Init(id string, model string) {
	l.Id = id
	l.EngineName = "Google"
	l.EngineUrl = "https://discoveryengine.googleapis.com/v1alpha/projects/764551478341/locations/global/collections/default_collection/dataStores/searchdemo_1687197269838/servingConfigs/default_search:search"
	l.EnginePrePrompt = ""
	l.EnginePostPrompt = ""
	l.EngineTimeout = 3500
	l.EngineToken = ""
	l.Callback = nil
	l.Model = model
}

func (l *GoogleES) LLMEngineName() string {
	return l.EngineName
}

func (l *GoogleES) LLMEngineUrl() string {
	return l.EngineUrl
}

func (l *GoogleES) LLMEnginePrePrompt() string {
	return l.EnginePrePrompt
}

func (l *GoogleES) LLMEnginePostPrompt() string {
	return l.EnginePostPrompt
}

func (l *GoogleES) LLMEngineModel() string {
	return l.Model
}

func (l *GoogleES) Completed() bool {
	return l.IsCompleted
}

func (l *GoogleES) Error() string {
	return l.DidError
}

func (l *GoogleES) LLMCitations() []string {
	return l.Citations
}

func (l *GoogleES) LLMResponse() string {
	if len(l.DidError) > 0 {
		return l.DidError
	} else {
		return l.QueryResponse
	}
}

func (l *GoogleES) LLMRuntime() int64 {
	return l.Runtime
}

func (l *GoogleES) Assign(query string, timeout int64, pre string, post string, callback WorkerCallback) {
	l.Callback = callback
	l.Query = query

	if timeout > 0 {
		l.EngineTimeout = timeout
	}

	// Do not allow pre/post on search
	/*
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

	*/
}

func (l *GoogleES) DoCallback() {
	isError := len(l.DidError) > 0
	if isError {
		l.Callback(l.Id, l.EngineName, l.DidError, l.Citations, isError, l.Runtime)
	} else {
		l.Callback(l.Id, l.EngineName, l.Response, l.Citations, isError, l.Runtime)
	}
}

func (l *GoogleES) Get_Google_Token() (string, error) {
	var token *oauth2.Token
	ctx := context.Background()
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
	credentials, err := auth.FindDefaultCredentials(ctx, scopes...)
	if err == nil {
		log.Printf("found default credentials. %v", credentials)
		token, err = credentials.TokenSource.Token()
		log.Printf("token access: %v", token.AccessToken)
		log.Printf("token: %v, err: %v", token, err)
		if err != nil {
			log.Print(err)
			return "", err
		}
		return token.AccessToken, nil
	} else {
		log.Printf("Error Authenticating: %v", err)
		return "", err
	}
}

func (l *GoogleES) Run(wg *sync.WaitGroup) {

	now := time.Now()
	defer func() {
		l.Runtime = time.Since(now).Milliseconds()
		l.DoCallback()
		wg.Done()
	}()

	token, err := l.Get_Google_Token()
	if err != nil {
		return
	}

	l.Request = ESRequest{}
	l.Request.Page_Size = 3
	l.Request.Offset = 0
	l.Request.Question = l.Query

	marshalled, err := json.Marshal(l.Request)
	if err != nil {
		l.DidError = fmt.Sprintf("impossible to marshall request: %s", err)
		log.Printf(l.DidError, err)
		return
	}

	req, err := http.NewRequest("POST", l.EngineUrl, bytes.NewReader(marshalled))
	if err != nil {
		l.DidError = fmt.Sprintf("impossible to build request: %s", err)
		log.Printf(l.DidError, err)
		return
	}

	// add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-functions-key", l.EngineToken)
	req.Header.Set("Authorization", "Bearer "+token)

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

	result := ESResponse{}

	// now we need to decode the string back to a response we can use
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		l.DidError = fmt.Sprintf("impossible to decode  response: %s", err)
		log.Printf(l.DidError, err)
		return
	}

	l.QueryResponse = ""
	l.Citations = []string{}
	for _, item := range result.Results {
		l.Citations = append(l.Citations, item.Id)
	}

}
