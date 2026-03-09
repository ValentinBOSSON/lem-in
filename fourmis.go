package main

import (
	"container/heap"
	"fmt"
	"os"
	"strings"
	"time"
)

func formatPath(path []*chamber) string {
	parts := make([]string, 0, len(path))
	for _, room := range path {
		parts = append(parts, room.ID)
	}
	return strings.Join(parts, " -> ")
}

func SimulateAntsOnPaths(ants []*ant, paths [][]*chamber, distribution []int, logPath string) error {
	if len(paths) == 0 {
		return fmt.Errorf("no paths provided")
	}

	logFile, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	// Assign ants to paths based on the distribution
	antIndex := 0
	for pathIdx := 0; pathIdx < len(paths); pathIdx++ {
		antsPerPath := distribution[pathIdx]
		for i := 0; i < antsPerPath && antIndex < len(ants); i++ {
			oneAnt := ants[antIndex]
			oneAnt.assignedPath = paths[pathIdx]
			oneAnt.pathPosition = 0
			oneAnt.currentRoom = paths[pathIdx][0]
			oneAnt.coordination = []int{paths[pathIdx][0].coordinates[0], paths[pathIdx][0].coordinates[1]}
			oneAnt.oldchambre = 0
			antIndex++
		}
	}

	// donne les conditions de fin et tout ce qui suit avec
	endRoom := paths[0][len(paths[0])-1]
	finished := 0
	turn := 1
	noProgressTurns := 0
	maxNoProgressTurns := len(ants) * 100

	// Heap pour guetter les rooms occupés et quand elles se libèrent
	roomOccupancy := make(OccupancyQueue, 0)
	heap.Init(&roomOccupancy)

	// Map pour verif si une room est occupé a un instant t
	roomFreeTime := make(map[string]int)

	for finished < len(ants) {
		moves := make([]string, 0)

		// Clean up expired occupancy entries from the heap
		for roomOccupancy.Len() > 0 {
			item := roomOccupancy[0]
			if item.timeFreed <= turn {
				heap.Pop(&roomOccupancy)
			} else {
				break
			}
		}

		// Try to move each ant
		for _, oneAnt := range ants {
			// Skip ants that have finished (reached the end)
			if oneAnt.currentRoom == endRoom {
				continue
			}

			// Skip ants without an assigned path
			if oneAnt.assignedPath == nil || len(oneAnt.assignedPath) < 2 {
				continue
			}

			// Try to move to next position on assigned path
			nextPos := oneAnt.pathPosition + 1
			if nextPos < len(oneAnt.assignedPath) {
				nextRoom := oneAnt.assignedPath[nextPos]

				// Check if next room is free or if it's the end room
				// Room is free if there's no occupancy entry or if the timeFreed has passed
				canMove := nextRoom == endRoom || roomFreeTime[nextRoom.ID] <= turn

				if canMove {
					oneAnt.currentRoom = nextRoom
					oneAnt.pathPosition = nextPos
					oneAnt.oldchambre++
					oneAnt.coordination = []int{nextRoom.coordinates[0], nextRoom.coordinates[1]}
					moves = append(moves, fmt.Sprintf("L%d-%s", oneAnt.ID, nextRoom.ID))

					// Mark room as occupied until next turn (but not the end room)
					if nextRoom != endRoom {
						occupancyItem := &RoomOccupancy{
							roomID:    nextRoom.ID,
							timeFreed: turn + 1,
							antID:     oneAnt.ID,
						}
						heap.Push(&roomOccupancy, occupancyItem)
						roomFreeTime[nextRoom.ID] = turn + 1
					}

					if nextRoom == endRoom {
						finished++
					}
				}
			}
		}

		// Track progress: increment counter each turn with no moves
		if len(moves) == 0 {
			noProgressTurns++
			if noProgressTurns > maxNoProgressTurns {
				return fmt.Errorf("simulation deadlocked: no progress after %d turns", maxNoProgressTurns)
			}
		} else {
			noProgressTurns = 0 // Reset counter when progress is made
		}

		// Always write a line, even if no moves (ants are waiting)
		var line string
		if len(moves) == 0 {
			line = fmt.Sprintf("%s | turn %d: (waiting)\n", time.Now().Format(time.RFC3339), turn)
		} else {
			line = fmt.Sprintf("%s | turn %d: %s\n", time.Now().Format(time.RFC3339), turn, strings.Join(moves, " "))
		}

		if _, err := logFile.WriteString(line); err != nil {
			return err
		}
		if err := logFile.Sync(); err != nil {
			return err
		}

		fmt.Print(line)
		turn++
	}

	return nil
}
func resetPathOccupancy(path []*chamber) {
	for index, room := range path {
		if index == 0 || index == len(path)-1 {
			continue
		}
		room.occupied = false
	}
}
