package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 2 || len(os.Args) < 2 {
		return
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading file %v", err)
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
	fmt.Println("\nAll chambers:")
	for id, chamber := range chambers {
		fmt.Printf("ID: %s, Position: %v, Connected to: ", id, chamber.coordinates)
		for _, conn := range chamber.tunnels {
			fmt.Printf("%s ", conn.ID)
		}
		fmt.Println()
	}
}
