package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func parseFile(file *os.File) (map[string]*chamber, int, *chamber, *chamber) {
	chambers := make(map[string]*chamber)
	temp := [][]string{}
	scanner := bufio.NewScanner(file)
	numAnts := 1
	var start, end *chamber
	nextIsStart := false
	nextIsEnd := false
	firstLine := true
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		if trimmedLine == "##start" {
			nextIsStart = true
			continue
		}
		if trimmedLine == "##end" {
			nextIsEnd = true
			continue
		}
		if strings.HasPrefix(trimmedLine, "#") && !strings.HasPrefix(trimmedLine, "##") {
			continue
		}
		if firstLine {
			numAnts, _ = strconv.Atoi(trimmedLine)
			firstLine = false
			continue
		}

		if strings.Contains(trimmedLine, "-") && len(strings.Fields(trimmedLine)) == 1 {
			foundTunnel := strings.Split(trimmedLine, "-")
			temp = append(temp, foundTunnel)
			continue
		}

		nodeParts := strings.Fields(trimmedLine)
		if len(nodeParts) == 3 {
			x, _ := strconv.Atoi(nodeParts[1])
			y, _ := strconv.Atoi(nodeParts[2])

			chamber := &chamber{
				ID:          nodeParts[0],
				coordinates: []int{x, y},
				occupied:    false,
				tunnels:     []*chamber{},
			}

			chambers[nodeParts[0]] = chamber

			if nextIsStart {
				start = chamber
				nextIsStart = false
			}
			if nextIsEnd {
				end = chamber
				nextIsEnd = false
			}
		}
	}

	for _, tunnel := range temp {
		node1 := tunnel[0]
		node2 := tunnel[1]
		chambers[node1].tunnels = append(chambers[node1].tunnels, chambers[node2])
		chambers[node2].tunnels = append(chambers[node2].tunnels, chambers[node1])
	}

	return chambers, numAnts, start, end
}

func createAnts(numAnts int, startChamber *chamber) []*ant {
	ants := make([]*ant, numAnts)
	for i := 0; i < numAnts; i++ {
		ants[i] = &ant{
			ID:           i + 1,
			coordination: []int{startChamber.coordinates[0], startChamber.coordinates[1]},
			oldchambre:   0,
		}
	}
	return ants
}
