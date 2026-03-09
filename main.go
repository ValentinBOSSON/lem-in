package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 || len(os.Args) < 2 {
		fmt.Println("ERROR: invalid data format")
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading file %v", err)
		fmt.Printf("\n")
		return
	}
	defer file.Close()
	chambers, numAnts, start, end := parseFile(file)
	ants := createAnts(numAnts, start)

	fmt.Printf("Found %d ants\n", numAnts)
	fmt.Printf("Start chamber: %s at %v\n", start.ID, start.coordinates)
	fmt.Printf("End chamber: %s at %v\n", end.ID, end.coordinates)
	if len(ants) > 0 {
		fmt.Printf("Created %d ants at %v\n", len(ants), ants[0].coordination)
	} else {
		fmt.Printf("Created 0 ants\n")
	}

	// Find vertex-disjoint paths for all ants
	paths := FindVertexDisjointPaths(start, end, numAnts)
	if len(paths) == 0 {
		fmt.Println("No paths found between start and end")
		return
	}
	fmt.Printf("Found %d disjoint paths\n", len(paths))
	for i, path := range paths {
		fmt.Printf("Path %d: ", i+1)
		for _, room := range path {
			fmt.Printf("%s ", room.ID)
		}
		fmt.Printf("(length: %d)\n", len(path))
	}

	// Calculate optimal ant distribution across paths
	distribution := OptimalAntDistribution(paths, numAnts)

	// Show the distribution
	totalExpectedTurns := 0
	for j := 0; j < len(paths); j++ {
		antsOnPath := distribution[j]
		expectedTurns := len(paths[j]) - 1 + antsOnPath - 1
		if expectedTurns > totalExpectedTurns {
			totalExpectedTurns = expectedTurns
		}
		fmt.Printf("Path %d: %d ants → expected completion in ~%d turns\n", j+1, antsOnPath, expectedTurns)
	}
	fmt.Printf("Expected total turns: %d\n", totalExpectedTurns)

	if err := SimulateAntsOnPaths(ants, paths, distribution, "deplacements.log"); err != nil {
		fmt.Printf("Simulation error: %v\n", err)
		return
	}
	fmt.Println("Log saved in deplacements.log")

	fmt.Println("\nAll chambers:")
	for id, chamber := range chambers {
		fmt.Printf("ID: %s, Position: %v, Connected to: ", id, chamber.coordinates)
		for _, conn := range chamber.tunnels {
			fmt.Printf("%s ", conn.ID)
		}
		fmt.Println()
	}
}
