package dashboard

import (
	"html/template"
	"net/http"
	"sync"
	"time"
)

// Dashboard -> Graph's -> line's -> point's.
type Dashboard struct {
	Title   string
	Refresh int // interval
	Graphs  []*Graph

	tpl   *template.Template
	mutex sync.Mutex
}

func New(title string, refreshInSecond int) *Dashboard {
	return &Dashboard{
		Title:   title,
		Refresh: refreshInSecond,
		Graphs:  make([]*Graph, 0),
		tpl:     template.Must(template.New("dashboard").Parse(tpl)),
	}
}

func (this *Dashboard) AddGraph(title string) *Graph {
	g := &Graph{
		Title: title,
		Lines: make([]*line, 0),
	}
	this.Graphs = append(this.Graphs, g)
	return g
}

func (this *Dashboard) Validate() error {
	for _, g := range this.Graphs {
		for _, l := range g.Lines {
			if l.dataSource == nil {
				return ErrEmptyDataSource
			}
		}
	}

	return nil
}

func (this *Dashboard) Launch(httpListenAddr string) error {
	if httpListenAddr == "" {
		return ErrEmptyHttpAddr
	}

	this.initHistory()

	http.HandleFunc("/", this.handleHttpRequest)
	return http.ListenAndServe(httpListenAddr, nil)
}

func (this *Dashboard) initHistory() {
	for _, g := range this.Graphs {
		for _, l := range g.Lines {
			if v, ok := l.dataSource.(WithHisotry); ok {
				for _, d := range v.History() {
					l.Points = append(l.Points, point{d[0], d[1]})
				}
			}
		}
	}

}

func (this *Dashboard) handleHttpRequest(w http.ResponseWriter, r *http.Request) {
	this.refreshData()
	this.tpl.Execute(w, this)
}

func (this *Dashboard) refreshData() {
	this.mutex.Lock()

	ts := int(time.Now().UnixNano() / 1e6)
	for _, g := range this.Graphs {
		for _, l := range g.Lines {
			l.Points = append(l.Points, point{ts, l.dataSource.Data()})
		}
	}

	this.mutex.Unlock()
}
