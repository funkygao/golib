package dashboard

type DataSource interface {
	Data() int
}

type point [2]int // [ts, data]

type line struct {
	Legend     string
	Points     []point
	dataSource DataSource
}

type Graph struct {
	Title string
	Lines []*line
}

func (this *Graph) AddLine(legend string, ds DataSource) {
	this.Lines = append(this.Lines,
		&line{dataSource: ds, Legend: legend, Points: make([]point, 0)})
}
