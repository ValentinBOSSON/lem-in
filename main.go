package main

<<<<<<< HEAD
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
	fmt.Printf("Created %d ants at %v\n", len(ants), ants[0].coordination)
	fmt.Println("\nAll chambers:")
	for id, chamber := range chambers {
		fmt.Printf("ID: %s, Position: %v, Connected to: ", id, chamber.coordinates)
		for _, conn := range chamber.tunnels {
			fmt.Printf("%s ", conn.ID)
		}
		fmt.Println()
=======
type Room struct {
	Ants     int
	Tunnels  [][]int
	Occupied bool
	Rooms    []string
	roompos  []int
}

func main() {
	ants := 1
	rooms := make(map[string]*Room)
	v := make(map[string][][]int)

	roomdata := map[string][]int{
		"A": {0, 3},
		"B": {2, 5},
		"C": {4, 0},
		"D": {8, 3},
	}

	antpos := [][]int{
		{0, 3},
	}

	v["A"] = [][]int{
		{0, 2},
		{2, 3},
		{3, 1},
	}

	for name, pos := range roomData {
		rooms[name] = &Room{
			Name:     name,
			Position: pos,
			Tunnels:  [][]int{}, // Will add connections later
			Occupied: isOccupied(pos, antpos),
			Ants:     0,
		}
>>>>>>> 29fe621 (fix:confict)
	}
}
