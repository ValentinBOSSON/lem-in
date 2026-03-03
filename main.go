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
	Name     string
	Position []int
}

func main() {
	rooms := make(map[string]*Room)
	v := make(map[string][][]int)

	roomdata := map[string][]int{
		"A": {0, 3},
		"B": {2, 5},
		"C": {4, 0},
		"D": {8, 3},
	}

	v["A"] = [][]int{
		{0, 2},
		{2, 3},
		{3, 1},
	}
	antpos := [][]int{
		{0, 3},
	}

	for name, pos := range roomdata {
		rooms[name] = &Room{
			Name:     name,
			Position: pos,
			Tunnels:  v[name],
			Occupied: IsOccupied(pos, antpos),
			Ants:     0,
		}
>>>>>>> 29fe621 (fix:confict)
	}

	for name, room := range rooms {
		println("Room:", name, "Position:", room.Position[0], room.Position[1], "Occupied:", room.Occupied)
	}

}
