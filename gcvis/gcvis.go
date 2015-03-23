package gcvis

import (
	"net/http"
)

var (
	g Graph
)

func Launch(addr string, title string) {
	if addr == "" {
		panic("addr can not be empty")
	}

	g = NewGraph(title, gcvis_tpl)

	http.HandleFunc("/", gcVisualize)
	go http.ListenAndServe(addr, nil)
}

func gcVisualize(w http.ResponseWriter, r *http.Request) {
	g.write(w)
}
