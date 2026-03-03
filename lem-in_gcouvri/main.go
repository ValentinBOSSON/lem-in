package main

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
	}
}
