package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func init(file *os.File) {
	chambers := make(map[string]*chamber)
	temp := [][]string{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}
		if strings.HasPrefix(trimmedLine, "#") {
			continue
		}
		if strings.Contains(trimmedLine, "-") {
			foundTunnel := strings.Split(trimmedLine, "-")
			temp = append(temp, foundTunnel)
		} else {
			nodeParts := strings.Fields(trimmedLine)
			if len(nodeParts) == 3 {
				x, _ := strconv.Atoi(nodeParts[1])
				y, _ := strconv.Atoi(nodeParts[2])
				chambers[nodeParts[0]] = &chamber{
					ID:          nodeParts[0],
					coordinates: []int{x, y},
					occupied:    false,
					tunnels:     []*chamber{},
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error scanning file %v", err)
	}
	for _, tunnel := range temp {
		node1 := tunnel[0]
		node2 := tunnel[1]
		chambers[node1].tunnels = append(chambers[node1].tunnels, chambers[node2])
		chambers[node2].tunnels = append(chambers[node2].tunnels, chambers[node1])
	}
}
