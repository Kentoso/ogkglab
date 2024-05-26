package sed

import "container/heap"

type Item struct {
	point    Point
	priority float32
	index    int
}

// PriorityQueue implements heap.Interface and holds Items
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

// NewPriorityQueue initializes a new priority queue
func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{}
	heap.Init(pq)
	return pq
}

func (pq *PriorityQueue) PushPoint(point Point, priority float32) {
	heap.Push(pq, &Item{
		point:    point,
		priority: priority,
	})
}

func (pq *PriorityQueue) PopPoint() (Point, float32) {
	item := heap.Pop(pq).(*Item)
	return item.point, item.priority
}

func (pq *PriorityQueue) IsEmpty() bool {
	return pq.Len() == 0
}
