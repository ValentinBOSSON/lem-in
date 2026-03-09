package main

import (
	"container/heap"
	"fmt"
	"sort"
)

type OccupancyQueue []*RoomOccupancy

func (oq OccupancyQueue) Len() int { return len(oq) }

func (oq OccupancyQueue) Less(i, j int) bool {
	return oq[i].timeFreed < oq[j].timeFreed
}

func (oq OccupancyQueue) Swap(i, j int) {
	oq[i], oq[j] = oq[j], oq[i]
	oq[i].index = i
	oq[j].index = j
}

func (oq *OccupancyQueue) Push(x any) {
	n := len(*oq)
	item := x.(*RoomOccupancy)
	item.index = n
	*oq = append(*oq, item)
}

func (oq *OccupancyQueue) Pop() any {
	old := *oq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*oq = old[0 : n-1]
	return item
}

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

func FindVertexDisjointPaths(start, end *chamber, maxPaths int) [][]*chamber {
	allShortestPaths := findAllShortestPaths(start, end)

	if len(allShortestPaths) == 0 {
		return nil
	}

	// Debug: print how many paths were found
	fmt.Printf("DEBUG: Found %d candidate paths\n", len(allShortestPaths))
	for i, p := range allShortestPaths {
		fmt.Printf("  Path %d: %s\n", i+1, formatPath(p))
	}

	// Sort paths by length (longer first) to try less-blocking paths first
	sort.Slice(allShortestPaths, func(i, j int) bool {
		return len(allShortestPaths[i]) > len(allShortestPaths[j])
	})

	// Try each shortest path as a starting point and see which gives the most disjoint paths
	var bestResult [][]*chamber

	for startIdx := 0; startIdx < len(allShortestPaths); startIdx++ {
		var result [][]*chamber
		blockedNodes := make(map[string]bool)

		// Start with this path
		result = append(result, allShortestPaths[startIdx])
		fmt.Printf("DEBUG: Starting with path %d: %s\n", startIdx+1, formatPath(allShortestPaths[startIdx]))

		for j := 1; j < len(allShortestPaths[startIdx])-1; j++ {
			blockedNodes[allShortestPaths[startIdx][j].ID] = true
		}

		// Iteratively find more paths
		for len(result) < maxPaths {
			path := FindFastestPathAvoiding(start, end, blockedNodes)
			fmt.Printf("DEBUG: After blocking %v, found path: %v\n", blockedNodes, path)

			if path == nil || len(path) < 2 {
				break
			}

			// Check for duplicate paths
			pathStr := formatPath(path)
			isDuplicate := false
			for _, existingPath := range result {
				if formatPath(existingPath) == pathStr {
					isDuplicate = true
					break
				}
			}

			if isDuplicate {
				break
			}

			result = append(result, path)

			// Block intermediate nodes
			for j := 1; j < len(path)-1; j++ {
				blockedNodes[path[j].ID] = true
			}
		}

		// Keep the combination that gives the most paths
		if len(result) > len(bestResult) {
			bestResult = result
		}
	}

	return bestResult
}

// Find all paths within a reasonable range of shortest distance using BFS
func findAllShortestPaths(start, end *chamber) [][]*chamber {
	// First, find the shortest distance
	dist := make(map[string]int)
	queue := make([]*chamber, 0)
	queue = append(queue, start)
	dist[start.ID] = 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		for _, next := range current.tunnels {
			if _, visited := dist[next.ID]; !visited {
				dist[next.ID] = dist[current.ID] + 1
				queue = append(queue, next)
			}
		}
	}

	shortestDist := dist[end.ID]
	maxDist := shortestDist + 2 // Include paths up to 2 hops longer

	// Now find paths within reasonable distance using DFS with cycle prevention
	var paths [][]*chamber
	var dfs func(*chamber, []*chamber, map[string]bool)
	dfs = func(current *chamber, path []*chamber, visited map[string]bool) {
		if current == end {
			newPath := make([]*chamber, len(path))
			copy(newPath, path)
			paths = append(paths, newPath)
			return
		}

		// Prune paths that are already too long
		if len(path)-1 >= maxDist {
			return
		}

		for _, next := range current.tunnels {
			// Only follow edges that keep us on a reasonable path and avoid cycles
			if !visited[next.ID] && dist[next.ID] <= dist[current.ID]+1 {
				newVisited := make(map[string]bool)
				for k, v := range visited {
					newVisited[k] = v
				}
				newVisited[next.ID] = true
				dfs(next, append(path, next), newVisited)
			}
		}
	}

	visited := make(map[string]bool)
	visited[start.ID] = true
	dfs(start, []*chamber{start}, visited)
	return paths
}

// Greedily find the best combination of vertex-disjoint paths
func findBestDisjointCombination(allPaths [][]*chamber, maxPaths int) [][]*chamber {
	var best [][]*chamber

	// Try starting with each path
	for startIdx := 0; startIdx < len(allPaths) && startIdx < maxPaths; startIdx++ {
		var result [][]*chamber
		usedNodes := make(map[string]bool)
		usedIndices := make(map[int]bool)

		// Add first path
		result = append(result, allPaths[startIdx])
		usedIndices[startIdx] = true
		for i := 1; i < len(allPaths[startIdx])-1; i++ {
			usedNodes[allPaths[startIdx][i].ID] = true
		}

		// Greedily add more paths, preferring shorter ones
		for len(result) < maxPaths {
			bestIdx := -1
			bestLen := int(^uint(0) >> 1) // max int

			for pathIdx := 0; pathIdx < len(allPaths); pathIdx++ {
				if usedIndices[pathIdx] {
					continue
				}

				// Check if this path is vertex-disjoint from already selected paths
				isDisjoint := true
				for i := 1; i < len(allPaths[pathIdx])-1; i++ {
					if usedNodes[allPaths[pathIdx][i].ID] {
						isDisjoint = false
						break
					}
				}

				// Prefer shorter paths first
				if isDisjoint && len(allPaths[pathIdx]) < bestLen {
					bestIdx = pathIdx
					bestLen = len(allPaths[pathIdx])
				}
			}

			if bestIdx == -1 {
				break
			}

			result = append(result, allPaths[bestIdx])
			usedIndices[bestIdx] = true
			for i := 1; i < len(allPaths[bestIdx])-1; i++ {
				usedNodes[allPaths[bestIdx][i].ID] = true
			}
		}

		// Keep the best result so far (prefer more paths, then shorter total length)
		if len(result) > len(best) {
			best = make([][]*chamber, len(result))
			copy(best, result)
		}
	}

	return best
}

// grosso modo la func a val
func FindPathWithConstraints(start, end *chamber, blockedEdges map[string]map[string]bool, blockedNodes map[string]bool) []*chamber {
	// heap l'arbre pour algo ddijkstra
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	cameFrom := make(map[string]*chamber)
	costSoFar := make(map[string]int)
	nodeToItem := make(map[string]*Item)

	// init du node en pushant l'index de base
	startItem := &Item{node: start, priority: 0}
	nodeToItem[start.ID] = startItem
	heap.Push(&pq, startItem)
	costSoFar[start.ID] = 0

	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*Item)
		current := item.node

		if item.priority > costSoFar[current.ID] {
			continue
		}

		if current.ID == end.ID {
			break
		}

		for _, next := range current.tunnels {

			// verif du blockage des edges et nodes
			if blockedEdges[current.ID] != nil && blockedEdges[current.ID][next.ID] {
				continue
			}
			if blockedNodes[next.ID] && next.ID != end.ID && next.ID != start.ID {
				continue
			}

			newCost := costSoFar[current.ID] + 1

			if cost, exists := costSoFar[next.ID]; !exists || newCost < cost {
				costSoFar[next.ID] = newCost
				cameFrom[next.ID] = current

				// si t'existe pas push si t'existe update la prio et fix l'arbre
				if !exists {
					item := &Item{node: next, priority: newCost}
					nodeToItem[next.ID] = item
					heap.Push(&pq, item)
				} else {
					existingItem, ok := nodeToItem[next.ID]
					if !ok {
						existingItem = &Item{node: next, priority: newCost}
						nodeToItem[next.ID] = existingItem
						heap.Push(&pq, existingItem)
					}
					existingItem.priority = newCost
					heap.Fix(&pq, existingItem.index)
				}
			}
		}
	}

	curr := end
	if _, exists := cameFrom[end.ID]; !exists {
		return nil
	}

	var path []*chamber
	for curr.ID != start.ID {
		path = append(path, curr)
		curr = cameFrom[curr.ID]
	}
	path = append(path, start)

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

// FindPathAvoidingEdges finds shortest path blocking specific edges
func FindPathAvoidingEdges(start, end *chamber, blockedEdges map[string]map[string]bool) []*chamber {
	return FindPathWithConstraints(start, end, blockedEdges, make(map[string]bool))
}

func FindFastestPathAvoiding(start, end *chamber, blockedNodes map[string]bool) []*chamber {
	return FindPathWithConstraints(start, end, make(map[string]map[string]bool), blockedNodes)
}

func FindFastestPath(start, end *chamber) []*chamber {
	return FindFastestPathAvoiding(start, end, make(map[string]bool))
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

// Uses a greedy approach: each ant is assigned to the path that will complete soonest
func OptimalAntDistribution(paths [][]*chamber, numAnts int) []int {
	distribution := make([]int, len(paths))

	// Greedy assignment: for each ant, assign to path with earliest completion
	for i := 0; i < numAnts; i++ {
		minCompletionTime := int(^uint(0) >> 1) // Max int
		bestPath := 0

		for j := 0; j < len(paths); j++ {
			// Completion time = path length + how many other ants are on this path
			// (they arrive serially at the rate of one per step)
			completionTime := len(paths[j]) - 1 + distribution[j]

			if completionTime < minCompletionTime {
				minCompletionTime = completionTime
				bestPath = j
			}
		}

		distribution[bestPath]++
	}

	return distribution
}
