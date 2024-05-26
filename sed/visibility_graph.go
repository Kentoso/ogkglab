package sed

type VisibilityGraph struct {
	AdjacencyMap map[Point][]Point
	S            Point
	T            Point
}

func NewVisibilityGraph(tr *TriangulationResult, s, t Point) VisibilityGraph {
	return VisibilityGraph{
		AdjacencyMap: tr.ToVertexAdjacencyMap(),
		S:            s,
		T:            t,
	}
}

func (vg *VisibilityGraph) ShortestEuclideanDistance() (map[Point]float32, []Point) {
	// Initialize the distance map and predecessor map
	distanceMap := make(map[Point]float32)
	predecessorMap := make(map[Point]Point)
	for v := range vg.AdjacencyMap {
		distanceMap[v] = float32(1<<31 - 1)
	}
	distanceMap[vg.S] = 0

	// Initialize the priority queue
	pq := NewPriorityQueue()
	pq.PushPoint(vg.S, 0)

	// ShortestEuclideanDistance's algorithm
	for !pq.IsEmpty() {
		v, _ := pq.PopPoint()

		// Relax the edges
		for _, u := range vg.AdjacencyMap[v] {
			if distanceMap[v]+v.Distance(u) < distanceMap[u] {
				distanceMap[u] = distanceMap[v] + v.Distance(u)
				predecessorMap[u] = v
				pq.PushPoint(u, distanceMap[u])
			}
		}
	}

	// Reconstruct the path from S to T
	path := []Point{}
	curr := vg.T

	// Traverse the path from T to S using the predecessor map
	for curr != vg.S {
		path = append([]Point{curr}, path...)
		curr = predecessorMap[curr]
	}
	path = append([]Point{vg.S}, path...)

	return distanceMap, path
}
