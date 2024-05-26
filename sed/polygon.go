package sed

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"math"
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

	for index, edge := range p.edges {
		line := canvas.NewLine(color.Black)
		line.Move(fyne.NewPos(edge.start.X, edge.start.Y))
		line.Resize(fyne.NewSize(edge.end.X-edge.start.X, edge.end.Y-edge.start.Y))
		objects = append(objects, line)

		arrowSize := 6.0

		angle := math.Atan2(float64(edge.end.Y-edge.start.Y), float64(edge.end.X-edge.start.X))
		arrowPoint1 := Point{
			X: edge.end.X - float32(arrowSize*math.Cos(angle-math.Pi/6)),
			Y: edge.end.Y - float32(arrowSize*math.Sin(angle-math.Pi/6)),
		}
		arrowPoint2 := Point{
			X: edge.end.X - float32(arrowSize*math.Cos(angle+math.Pi/6)),
			Y: edge.end.Y - float32(arrowSize*math.Sin(angle+math.Pi/6)),
		}

		arrowLine1 := canvas.NewLine(color.Black)
		arrowLine1.Position1 = fyne.NewPos(edge.end.X, edge.end.Y)
		arrowLine1.Position2 = fyne.NewPos(arrowPoint1.X, arrowPoint1.Y)
		objects = append(objects, arrowLine1)

		arrowLine2 := canvas.NewLine(color.Black)
		arrowLine2.Position1 = fyne.NewPos(edge.end.X, edge.end.Y)
		arrowLine2.Position2 = fyne.NewPos(arrowPoint2.X, arrowPoint2.Y)
		objects = append(objects, arrowLine2)

		midX := (edge.start.X + edge.end.X) / 2
		midY := (edge.start.Y + edge.end.Y) / 2

		text := canvas.NewText(fmt.Sprintf("%d", index), color.Black)
		text.TextSize = 12
		text.Move(fyne.NewPos(midX-6, midY-6))
		objects = append(objects, text)
	}

	for index, vertex := range p.vertices {
		circle := canvas.NewCircle(color.Black)
		circle.Resize(fyne.NewSize(5, 5))
		circle.Move(fyne.NewPos(vertex.X-2.5, vertex.Y-2.5))
		objects = append(objects, circle)

		text := canvas.NewText(fmt.Sprintf("%d", index), color.RGBA{0, 255, 0, 255})
		text.TextSize = 12
		text.Move(fyne.NewPos(vertex.X-3, vertex.Y-3))
		objects = append(objects, text)
	}

	return container.NewWithoutLayout(objects...)
}
