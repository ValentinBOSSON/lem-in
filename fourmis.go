package main

import (
	"fmt"
	"os"
	"sort"
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

func resetPathOccupancy(path []*chamber) {
	for index, room := range path {
		if index == 0 || index == len(path)-1 {
			continue
		}
		room.occupied = false
	}
}

func SimulateAntsOnPath(ants []*ant, path []*chamber, logPath string) error {
	if len(path) < 2 {
		return fmt.Errorf("invalid path")
	}

	startRoom := path[0]
	endRoom := path[len(path)-1]

	logFile, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	resetPathOccupancy(path)
	for _, oneAnt := range ants {
		oneAnt.oldchambre = 0
		oneAnt.currentRoom = startRoom
		oneAnt.coordination = []int{path[0].coordinates[0], path[0].coordinates[1]}
	}

	finished := 0
	turn := 1

	for finished < len(ants) {
		endRoom.occupied = false

		sort.SliceStable(ants, func(i, j int) bool {
			if ants[i].oldchambre == ants[j].oldchambre {
				return ants[i].ID < ants[j].ID
			}
			return ants[i].oldchambre > ants[j].oldchambre
		})

		moves := make([]string, 0)

		for _, oneAnt := range ants {
			if oneAnt.currentRoom == endRoom {
				continue
			}

			bestPath := FindFastestPath(oneAnt.currentRoom, endRoom)
			if len(bestPath) < 2 {
				continue
			}

			nextRoom := bestPath[1]
			if nextRoom.occupied && nextRoom != endRoom {
				continue
			}

			if oneAnt.currentRoom != startRoom && oneAnt.currentRoom != endRoom {
				oneAnt.currentRoom.occupied = false
			}
			if nextRoom != startRoom && nextRoom != endRoom {
				nextRoom.occupied = true
			}

			oneAnt.currentRoom = nextRoom
			oneAnt.oldchambre++
			oneAnt.coordination = []int{nextRoom.coordinates[0], nextRoom.coordinates[1]}
			moves = append(moves, fmt.Sprintf("L%d-%s", oneAnt.ID, nextRoom.ID))

			if nextRoom == endRoom {
				endRoom.occupied = true
				finished++
			}
		}

		if len(moves) == 0 {
			return fmt.Errorf("simulation blocked: no ant could move")
		}

		line := fmt.Sprintf("%s | turn %d: %s\n", time.Now().Format(time.RFC3339), turn, strings.Join(moves, " "))
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
