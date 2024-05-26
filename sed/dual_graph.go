package sed

import "fmt"

type DualGraph struct {
	adjacencyList map[int]map[int]struct{} // Adjacency list representation of the dual graph using sets
}

func NewDualGraph() DualGraph {
	return DualGraph{
		adjacencyList: make(map[int]map[int]struct{}),
	}
}

func BuildGraphDual(tr TriangulationResult) DualGraph {
	dg := NewDualGraph()

	for i, triangle := range tr.triangles {
		for _, neighbor := range triangle.neighbors {
			if neighbor != -1 {
				dg.AddEdge(i, neighbor)
			}
		}
	}

	return dg
}

func (dg *DualGraph) AddEdge(t1, t2 int) {
	if dg.adjacencyList[t1] == nil {
		dg.adjacencyList[t1] = make(map[int]struct{})
	}
	if dg.adjacencyList[t2] == nil {
		dg.adjacencyList[t2] = make(map[int]struct{})
	}
	dg.adjacencyList[t1][t2] = struct{}{}
	dg.adjacencyList[t2][t1] = struct{}{}
}

func (dg *DualGraph) RemoveEdge(t1, t2 int) {
	delete(dg.adjacencyList[t1], t2)
	delete(dg.adjacencyList[t2], t1)
	if len(dg.adjacencyList[t1]) == 0 {
		delete(dg.adjacencyList, t1)
	}
	if len(dg.adjacencyList[t2]) == 0 {
		delete(dg.adjacencyList, t2)
	}
}

func (dg *DualGraph) Copy() DualGraph {
	newGraph := NewDualGraph()
	for node, neighbors := range dg.adjacencyList {
		newGraph.adjacencyList[node] = make(map[int]struct{})
		for neighbor := range neighbors {
			newGraph.adjacencyList[node][neighbor] = struct{}{}
		}
	}
	return newGraph
}

func (dg *DualGraph) Simplify() DualGraph {
	removed := true
	withoutFirstDegree := dg.Copy()
	for removed {
		removed = false
		for node, neighbors := range withoutFirstDegree.adjacencyList {
			if len(neighbors) == 1 {
				// Remove degree-1 node
				for neighbor := range neighbors {
					withoutFirstDegree.RemoveEdge(node, neighbor)
				}
				removed = true
			}
		}
	}

	withoutSecondDegree := withoutFirstDegree.Copy()
	secondDegreeNodes := []int{}
	for node, neighbors := range withoutSecondDegree.adjacencyList {
		if len(neighbors) == 2 {
			secondDegreeNodes = append(secondDegreeNodes, node)
		}
	}

	for _, node := range secondDegreeNodes {
		neighbors := withoutSecondDegree.adjacencyList[node]
		if len(neighbors) == 2 {
			// Replace degree-2 node with a single edge between its neighbors
			fmt.Printf("Removing degree-2 node %d\n", node)
			var neighbor1, neighbor2 int
			count := 0
			for neighbor := range neighbors {
				if count == 0 {
					neighbor1 = neighbor
				} else {
					neighbor2 = neighbor
				}
				count++
			}
			fmt.Printf("Removing edges %d-%d and %d-%d\n", node, neighbor1, node, neighbor2)
			fmt.Printf("Adding edge %d-%d\n", neighbor1, neighbor2)

			withoutSecondDegree.RemoveEdge(node, neighbor1)
			withoutSecondDegree.RemoveEdge(node, neighbor2)
			withoutSecondDegree.AddEdge(neighbor1, neighbor2)
		}
	}

	return withoutSecondDegree
}
