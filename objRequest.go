package main

type objRequest struct {
	User    string   `json:"user"`
	Query   string   `json:"query"`
	Engines []string `json:"engines"`
	Debug   bool     `json:"debug"`
	Pre     string   `json:"pre"`
	Post    string   `json:"post"`
}
