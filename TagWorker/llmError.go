package tagworker

import (
	"strings"
	"sync"
	"time"
)

type LLMEngineError struct {
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
}

func (l *LLMEngineError) Init(id string, model string) {
	l.Id = id
	l.EngineName = "LLMError"
	l.EngineUrl = ""
	l.EnginePrePrompt = ""
	l.EnginePostPrompt = ""
	l.EngineTimeout = 3500
	l.EngineToken = ""
	l.Callback = nil
	l.Model = model
}

func (l *LLMEngineError) LLMEngineName() string {
	return l.EngineName
}

func (l *LLMEngineError) LLMEngineUrl() string {
	return l.EngineUrl
}

func (l *LLMEngineError) LLMEnginePrePrompt() string {
	return l.EnginePrePrompt
}

func (l *LLMEngineError) LLMEnginePostPrompt() string {
	return l.EnginePostPrompt
}

func (l *LLMEngineError) LLMEngineModel() string {
	return l.Model
}

func (l *LLMEngineError) Completed() bool {
	return l.IsCompleted
}

func (l *LLMEngineError) Error() string {
	return l.DidError
}

func (l *LLMEngineError) LLMCitations() []string {
	return l.Citations
}

func (l *LLMEngineError) LLMResponse() string {
	if len(l.DidError) > 0 {
		return l.DidError
	} else {
		return l.QueryResponse
	}
}

func (l *LLMEngineError) LLMRuntime() int64 {
	return l.Runtime
}

func (l *LLMEngineError) Assign(query string, timeout int64, pre string, post string, callback WorkerCallback) {
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

func (l *LLMEngineError) Run(wg *sync.WaitGroup) {

	now := time.Now()
	defer func() {
		l.Runtime = time.Since(now).Milliseconds()
		wg.Done()
	}()
}
