package tagworker

import "sync"

type LLMEngine interface {
	Init(id string, model string)
	LLMEngineName() string
	LLMEngineUrl() string
	LLMEngineModel() string
	LLMEnginePrePrompt() string
	LLMEnginePostPrompt() string

	LLMCitations() []string
	LLMResponse() string
	LLMRuntime() int64

	Assign(query string, timeout int64, pre string, post string, callback WorkerCallback)
	Run(wg *sync.WaitGroup)
	Completed() bool
	Error() string
}
