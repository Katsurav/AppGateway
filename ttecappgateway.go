package main

import (
	"TTECAppGateway/KIALogging"
	"TTECAppGateway/TagWorker"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	guid "github.com/google/uuid"
	router "github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	auth "golang.org/x/oauth2/google"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	//go:embed resources
	res   embed.FS
	pages = map[string]string{
		"basic":               "resources/basic.html",
		"all":                 "resources/all.html",
		"prompt":              "resources/prompt.html",
		"promptall":           "resources/promptall.html",
		"script.js":           "resources/script.js",
		"style.css":           "resources/style.css",
		"query":               "resources/query.html",
		"viking.png":          "resources/viking.png",
		"bg.png":              "resources/bg.png",
	}

	DBCloud = "/cloudsql/runningtortoiselab:us-central1:kiaalpha_"
	UseDB   = false
	DB      = KIALogging.KIALogger{}
)

func Google_Token() {
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
		}
	} else {
		log.Printf("Error Authenticating: %v", err)
	}
}

// Basic Auth

func BasicAuth(h router.Handle, requiredUser, requiredPassword string) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps router.Params) {
		// Get the Basic Authentication credentials
		user, password, hasAuth := r.BasicAuth()

		if hasAuth && user == requiredUser && password == requiredPassword {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func ShowToken(w http.ResponseWriter, r *http.Request, _ router.Params) {
	fmt.Fprint(w, "Token!\n")
	Google_Token()
}

func Web_Serve(w http.ResponseWriter, r *http.Request, p router.Params) {
	name := p.ByName("Name")

	var binary = false
	if strings.Contains(name, ".css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.Contains(name, ".png") {
		w.Header().Set("Content-Type", "image/png")
		log.Printf("Binary: true")
		binary = true
	} else if strings.Contains(name, ".js") {
		w.Header().Set("Content-Type", "text/javascript")
		log.Printf("Javascript: true")
		binary = true
	} else {
		w.Header().Set("Content-Type", "text/html")
	}

	page, ok := pages[name]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if binary {
		data, _ := res.ReadFile(page)
		log.Printf("Read File: %s, Size: %d", page, len(data))
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else {

		tpl, err := template.ParseFS(res, page)
		if err != nil {
			log.Printf("page %s not found in pages cache...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		data := map[string]interface{}{
			"userAgent": r.UserAgent(),
		}
		if err := tpl.Execute(w, data); err != nil {
			return
		}
	}

}

func Index(w http.ResponseWriter, r *http.Request, _ router.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func Health(w http.ResponseWriter, r *http.Request, _ router.Params) {
	fmt.Fprint(w, "Healthy!\n")
}

func Config(w http.ResponseWriter, r *http.Request, _ router.Params) {
	response := "<!DOCTYPE html>\n<html lang=\"en\">\n<body>"

	var ii tagworker.LLMEngine

	ii = &tagworker.AzureLLM{}
	ii.Init("", "gpt-35-turbo")
	response = response + MakeConfigLLM(ii)

	ii = &tagworker.AzureLLM{}
	ii.Init("", "gpt-35-turbo-preview")
	response = response + MakeConfigLLM(ii)

	ii = &tagworker.AzureLLM{}
	ii.Init("", "text-bison")
	response = response + MakeConfigLLM(ii)

	ii = &tagworker.GooglePalm{}
	ii.Init("", "palm")
	response = response + MakeConfigLLM(ii)

	ii = &tagworker.GoogleES{}
	ii.Init("", "es")
	response = response + MakeConfigLLM(ii)

	response = response + "</body> \n </html> \n"

	bb := []byte(response)
	w.Write(bb)
}

func MakeConfigLLM(e tagworker.LLMEngine) string {
	response := ""
	response = "<span style=\"white-space: pre;\">"
	response = response + "Name : " + e.LLMEngineName() + "\n"
	response = response + "Model: " + e.LLMEngineModel() + "\n"
	response = response + "Pre  : " + e.LLMEnginePrePrompt() + "\n"
	response = response + "Post : " + e.LLMEnginePostPrompt() + "\n"
	response = response + "Url  : " + e.LLMEngineUrl() + "\n"
	response = response + "\n\n</span>"
	return response
}

func Hello(w http.ResponseWriter, r *http.Request, ps router.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {

	if UseDB {
		// DB
		passed := DB.Connect(DBCloud)
		UseDB = passed
	}

	Router := router.New()

	Router.HandleOPTIONS = true
	Router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
			header.Add("Access-Control-Allow-Origin", "*")
			header.Add("Access-Control-Max-Age", "3600")
			header.Add("Access-Control-Allow-Headers", "access-control-allow-origin, x-requested-with, content-type")
			log.Println("OPTIONS Requested")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	Router.GET("/", Index)
	Router.GET("/health", Health)
	Router.GET("/token", ShowToken)
	Router.GET("/config", Config)
	Router.POST("/query", Query)

	Router.GET("/hello/:name", Hello)

	Router.GET("/web/:Name", Web_Serve)

	// router.GET("/protected/", BasicAuth(Protected, user, pass))

	port := 8080

	log.Printf("Open http://localhost:%d in the browser", port)

	log.Fatal(http.ListenAndServe(":8080", Router))

}

func Query(w http.ResponseWriter, r *http.Request, _ router.Params) {
	header := w.Header()
	header.Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Max-Age", "3600")
	header.Add("Access-Control-Allow-Headers", "access-control-allow-origin, x-requested-with, content-type")

	reqBody, _ := io.ReadAll(r.Body)
	log.Println(string(reqBody))
	var post objRequest
	err := json.Unmarshal(reqBody, &post)
	if err != nil {
		// We errored
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
	} else {
		worker := &tagworker.Worker{}
		uuidWithHyphen := guid.New()
		id := uuidWithHyphen.String()

		if UseDB {
			DB.QueryCreate(id, "", post.Query, post.Pre, post.Post)
		}

		err := worker.Init(id, post.Engines)
		if err != nil {
			// We errored
			w.WriteHeader(503)
			w.Write([]byte(err.Error()))
		} else {
			worker.Assign(post.Query, 5000, post.Pre, post.Post)
			worker.Run()

			Response := objResponse{}
			Response.RequestId = worker.Id
			Response.Query = post.Query
			Response.Results = []objResult{}
			rank := 0
			for key, element := range worker.Workers {
				nr := objResult{}
				rank = rank + 1
				nr.LLMID = key
				nr.LLMName = element.LLMEngineName() + "/" + element.LLMEngineModel()
				nr.DidError = len(element.Error()) > 0
				nr.LLMCitations = element.LLMCitations()
				nr.LLMResponse = element.LLMResponse()
				nr.LLMCompute = element.LLMRuntime()
				Response.Results = append(Response.Results, nr)

				if UseDB {
					DB.ResponseCreate(key, rank, element.LLMResponse(), element.LLMEngineModel(), element.LLMRuntime())
					cit := 0
					for _, citation := range element.LLMCitations() {
						cit = cit + 1
						DB.CitationCreate(key, cit, citation)
					}
				}
			}

			marshalled, err := json.Marshal(Response)
			if err != nil {
				w.WriteHeader(503)
				w.Write([]byte(err.Error()))
			}

			header.Add("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)

			w.Write(marshalled)
		}
	}

}
