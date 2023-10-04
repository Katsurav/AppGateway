package tagworker

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type WorkerCallback func(Id string, Engine string, Response string, Citations []string, DidError bool, Runtime int64)

type Worker struct {
	Engines []string
	Workers map[string]LLMEngine
	Id      string
}

func (w *Worker) Init(id string, engines []string) error {

	w.Id = id
	w.Workers = make(map[string]LLMEngine)

	if engines == nil || len(engines) == 0 {
		return errors.New("no LLM Engines Provided")
	}

	var istring string
	for index, element := range engines {
		// index is the index where we are
		// element is the element from someSlice for where we are
		istring = id + "." + strconv.Itoa(index)
		switch element {
		case "gpt-35-turbo":
			w.Workers[istring] = &AzureLLM{}
			w.Workers[istring].Init(istring, "gpt-35-turbo")
		case "gpt-35-turbo-preview":
			w.Workers[istring] = &AzureLLM{}
			w.Workers[istring].Init(istring, "gpt-35-turbo-preview")
		case "text-bison":
			w.Workers[istring] = &GooglePalm{}
			w.Workers[istring].Init(istring, "text-bison")
		case "palm":
			w.Workers[istring] = &GooglePalm{}
			w.Workers[istring].Init(istring, "palm")
		case "es":
			w.Workers[istring] = &GoogleES{}
			w.Workers[istring].Init(istring, "es")

		default:
			ee := LLMEngineError{}
			ee.DidError = "unknown LLM Engine Requested + [" + element + "]"
			ee.Model = element
			w.Workers[istring] = &ee
			w.Workers[istring].Init(istring, "")
		}

	}

	return nil
}

func (w *Worker) Assign(query string, timeout int64, pre string, post string) {
	for key, element := range w.Workers {
		fmt.Println("Assign", "", "Key:", key, "=>", "Element:", element)
		element.Assign(query, timeout, pre, post, w.EngineCallback)
	}
}

func (w *Worker) Run() {
	// Need to Convert this to GroupWait and run all 3 in parallel functions
	wg := new(sync.WaitGroup)

	for key, element := range w.Workers {
		wg.Add(1)
		fmt.Println("Run", "", "Key:", key, "=>", "Element:", element)
		go element.Run(wg)
	}

	wg.Wait()

	fmt.Println("Completed Execution")
}

func (w *Worker) EngineCallback(Id string, Engine string, Response string, Citations []string, DidError bool, Runtime int64) {
	fmt.Printf("Callback %s %s %d\n", Engine, Response, Runtime)
}
