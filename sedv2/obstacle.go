package sedv2

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

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

	return Obstacle{Vertices: points}
}

func (o Obstacle) Translate(x float32, y float32) Obstacle {
	for i := range o.Vertices {
		o.Vertices[i].X += x
		o.Vertices[i].Y += y
	}
	return o
}
