package main

import (
	"container/heap"
)

type Item struct {
	node     *chamber
	priority int
	index    int
}

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

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func FindFastestPath(start, end *chamber) []*chamber {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	cameFrom := make(map[string]*chamber)
	costSoFar := make(map[string]int)

	heap.Push(&pq, &Item{node: start, priority: 0})
	costSoFar[start.ID] = 0

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*Item).node

		if current.ID == end.ID {
			break
		}

		for _, next := range current.tunnels {
			// Skip occupied rooms to allow finding vertex-disjoint paths progressively
			if next.occupied && next.ID != end.ID && next.ID != start.ID {
				continue
			}

			newCost := costSoFar[current.ID] + 1
			if cost, exists := costSoFar[next.ID]; !exists || newCost < cost {
				costSoFar[next.ID] = newCost
				cameFrom[next.ID] = current
				heap.Push(&pq, &Item{node: next, priority: newCost})
			}
		}
	}

	curr := end
	if _, exists := cameFrom[end.ID]; !exists {
		return nil // No path found
	}

	var path []*chamber
	for curr.ID != start.ID {
		path = append(path, curr)
		curr = cameFrom[curr.ID]
	}
	path = append(path, start)

	// Reverse path so it goes from start -> end
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func IsOccupied(x int, y int, antpos *ant) bool {
	if antpos.coordination[0] == x && antpos.coordination[1] == y {
		return true
	}
	return false
}

func Canmove(room *chamber, antpos [][]int) bool {
	if room.occupied == false {
		return true
	}
	return false
}

/*func OptimizedPath(end *chamber) {
	if chamber.occupied == true {

	}
}
*/
