package main

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
	}

	for name, room := range rooms {
		println("Room:", name, "Position:", room.Position[0], room.Position[1], "Occupied:", room.Occupied)
	}

}
