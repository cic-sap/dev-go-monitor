package main

import (
	"fmt"
	pluginsimple "github.com/cic-sap/dev-go-monitor/plugin-simple"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})
	mux.HandleFunc("/user/:id", func(w http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(w, "get user :%s\n", request.URL)
	})
	h := pluginsimple.Init(mux)
	http.ListenAndServe(":8290", h)
}
