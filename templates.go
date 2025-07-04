package main

type ConfigJSON struct {
	Handlers []Handlers `json:"handlers"`
}

type Handlers struct {
	Title    string     `json:"title"`
	Method   string     `json:"method"`
	Path     string     `json:"path"`
	Examples []Examples `json:"examples"`
}

type Examples struct {
	Header     interface{} `json:"header"`
	Body       interface{} `json:"body"`
	Wanted     interface{} `json:"wanted"`
	WantedCode int         `json:"wanted_code"`
}

type ApifoxJSON struct {
}
