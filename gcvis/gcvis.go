package gcvis

import (
	"net/http"
)

var (
	g *graph
)

// Launch will start GC visualization on specifed http addr.
// title is the HTML title while refresh is the second interval
// between auto refresh.
func Launch(httpAddr, title string, refresh int) {
	if httpAddr == "" {
		panic("httpAddr can not be empty")
	}

	g = newGraph(title, gcvis_tpl, refresh)

	http.HandleFunc("/", gcVisualize)
	go http.ListenAndServe(httpAddr, nil)
}

func gcVisualize(w http.ResponseWriter, r *http.Request) {
	g.write(w)
}
