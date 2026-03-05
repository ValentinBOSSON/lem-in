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

	logFile, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer logFile.Close()

	resetPathOccupancy(path)
	for _, oneAnt := range ants {
		oneAnt.oldchambre = 0
		oneAnt.coordination = []int{path[0].coordinates[0], path[0].coordinates[1]}
	}

	finished := 0
	turn := 1

	for finished < len(ants) {
		sort.SliceStable(ants, func(i, j int) bool {
			if ants[i].oldchambre == ants[j].oldchambre {
				return ants[i].ID < ants[j].ID
			}
			return ants[i].oldchambre > ants[j].oldchambre
		})

		moves := make([]string, 0)

		for _, oneAnt := range ants {
			if oneAnt.oldchambre >= len(path)-1 {
				continue
			}

			currentIndex := oneAnt.oldchambre
			nextIndex := currentIndex + 1
			nextRoom := path[nextIndex]

			canMove := nextIndex == len(path)-1 || !nextRoom.occupied
			if !canMove {
				continue
			}

			if currentIndex > 0 && currentIndex < len(path)-1 {
				path[currentIndex].occupied = false
			}
			if nextIndex > 0 && nextIndex < len(path)-1 {
				nextRoom.occupied = true
			}

			oneAnt.oldchambre = nextIndex
			oneAnt.coordination = []int{nextRoom.coordinates[0], nextRoom.coordinates[1]}
			moves = append(moves, fmt.Sprintf("L%d-%s", oneAnt.ID, nextRoom.ID))

			if nextIndex == len(path)-1 {
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
