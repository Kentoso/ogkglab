package sed

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"image/color"
	"math"
	"math/rand"
	"slices"
	"sort"
	"time"
)

type Obstacle struct {
	vertices []Point
}

func CreateRandomObstacle(numPoints int, minX, minY, maxX, maxY float32) Obstacle {
	rand.Seed(time.Now().UnixNano())

	points := make([]Point, numPoints)
	for i := 0; i < numPoints; i++ {
		points[i] = Point{
			X: minX + rand.Float32()*(maxX-minX),
			Y: minY + rand.Float32()*(maxY-minY),
		}
	}

	var centroid Point
	for _, p := range points {
		centroid.X += p.X
		centroid.Y += p.Y
	}
	centroid.X /= float32(numPoints)
	centroid.Y /= float32(numPoints)

	sort.Slice(points, func(i, j int) bool {
		angle1 := math.Atan2(float64(points[i].Y-centroid.Y), float64(points[i].X-centroid.X))
		angle2 := math.Atan2(float64(points[j].Y-centroid.Y), float64(points[j].X-centroid.X))
		return angle1 < angle2
	})

	return Obstacle{vertices: points}
}

type Map struct {
	obstacles []Obstacle
	s         Point
	t         Point
}

func NewMap(s, t Point) *Map {
	return &Map{nil, s, t}
}

func (m *Map) AddObstacle(obstacle Obstacle) {
	m.obstacles = append(m.obstacles, obstacle)
}

func (m *Map) ToPolygon() *Polygon {
	var vertices []Point
	var edges []Segment

	for _, obstacle := range m.obstacles {
		for i := 0; i < len(obstacle.vertices); i++ {
			vertices = append(vertices, obstacle.vertices[i])
			if i > 0 {
				edges = append(edges, Segment{start: obstacle.vertices[i-1], end: obstacle.vertices[i]})
			}
		}
		edges = append(edges, Segment{start: obstacle.vertices[len(obstacle.vertices)-1], end: obstacle.vertices[0]})
	}

	//vertices = append(vertices, m.s, m.t)

	minX, minY := vertices[0].X, vertices[0].Y
	maxX, maxY := vertices[0].X, vertices[0].Y

	for _, v := range append(vertices, []Point{m.s, m.t}...) {
		if v.X < minX {
			minX = v.X
		}
		if v.X > maxX {
			maxX = v.X
		}
		if v.Y < minY {
			minY = v.Y
		}
		if v.Y > maxY {
			maxY = v.Y
		}
	}

	var margin float32 = 10.0
	minX -= margin
	maxX += margin
	minY -= margin
	maxY += margin

	boundingBoxVertices := []Point{
		{X: minX, Y: minY},
		{X: maxX, Y: minY},
		{X: maxX, Y: maxY},
		{X: minX, Y: maxY},
	}

	vertices = append(vertices, boundingBoxVertices...)

	for i := 0; i < len(boundingBoxVertices); i++ {
		edges = append(edges, Segment{
			start: boundingBoxVertices[i],
			end:   boundingBoxVertices[(i+1)%len(boundingBoxVertices)],
		})
	}

	edges = m.addBridges(vertices[len(vertices)-4:], edges)

	return &Polygon{
		vertices: vertices,
		edges:    edges,
	}
}

func (m *Map) addBridges(vertices []Point, edges []Segment) []Segment {
	// Create bridges for each hole

	currentEdges := slices.Clone(edges)
	for i := 0; i < len(m.obstacles); i++ {
		hole := m.obstacles[i].vertices
		//polygon := m.obstacles[0].vertices
		minDistance := math.MaxFloat64
		var bestBridge Segment

		verticesToCheck := slices.Clone(vertices)
		for j := 0; j < len(m.obstacles); j++ {
			if i == j {
				continue
			}

			verticesToCheck = append(verticesToCheck, m.obstacles[j].vertices...)
		}

		for _, p := range verticesToCheck {
			for _, h := range hole {
				bridge := Segment{start: p, end: h}
				valid := true
				distance := math.Hypot(float64(bridge.end.X-bridge.start.X), float64(bridge.end.Y-bridge.start.Y))

				if distance == 0 {
					continue
				}
				// Check if this bridge intersects any existing edges
				for _, e := range currentEdges {
					if segmentsIntersect(bridge.start, bridge.end, e.start, e.end) {
						valid = false
						break
					}
				}

				if valid {
					if distance > 0 && distance < minDistance {
						minDistance = distance
						bestBridge = bridge
					}
				}
			}
		}

		// Add the best bridge to the list of edges
		//bestBridgeReversed := Segment{start: bestBridge.end, end: bestBridge.start}

		edges = append(edges, bestBridge)
		//edges = append(edges, bestBridgeReversed)

		fmt.Println("Added ", bestBridge)
	}

	return edges
}

func (m *Map) Draw() fyne.CanvasObject {
	objects := []fyne.CanvasObject{}

	// Draw obstacles
	for _, obstacle := range m.obstacles {
		for i := 0; i < len(obstacle.vertices); i++ {
			start := obstacle.vertices[i]
			end := obstacle.vertices[(i+1)%len(obstacle.vertices)]
			line := canvas.NewLine(color.Black)
			line.Position1 = fyne.NewPos(start.X, start.Y)
			line.Position2 = fyne.NewPos(end.X, end.Y)
			objects = append(objects, line)
		}
	}

	// Draw start and end points
	points := []struct {
		point Point
		label string
	}{
		{m.s, "S"},
		{m.t, "T"},
	}

	for _, p := range points {
		// Draw the circle
		circle := canvas.NewCircle(color.RGBA{255, 0, 0, 255})
		circle.Resize(fyne.NewSize(5, 5))
		circle.Move(fyne.NewPos(p.point.X-2.5, p.point.Y-2.5))
		objects = append(objects, circle)

		// Draw the label
		text := canvas.NewText(p.label, color.Black)
		text.TextSize = 12
		text.Move(fyne.NewPos(p.point.X+5, p.point.Y-6)) // Positioning the text near the point
		objects = append(objects, text)
	}

	return container.NewWithoutLayout(objects...)
}
