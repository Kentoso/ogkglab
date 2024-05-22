package sed

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
)

type Segment struct {
	start, end Point
}

type Polygon struct {
	vertices []Point
	edges    []Segment
}

func (p *Polygon) Draw() fyne.CanvasObject {
	objects := []fyne.CanvasObject{}

	for _, edge := range p.edges {
		line := canvas.NewLine(color.Black)
		line.Move(fyne.NewPos(edge.start.X, edge.start.Y))
		line.Resize(fyne.NewSize(edge.end.X-edge.start.X, edge.end.Y-edge.start.Y))
		objects = append(objects, line)
	}

	for _, vertex := range p.vertices {
		circle := canvas.NewCircle(color.Black)
		circle.Resize(fyne.NewSize(5, 5))
		circle.Move(fyne.NewPos(vertex.X-2.5, vertex.Y-2.5))
		objects = append(objects, circle)
	}

	return container.NewWithoutLayout(objects...)
}
