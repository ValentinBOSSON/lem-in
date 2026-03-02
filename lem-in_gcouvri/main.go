package main


type Room struct {
	Ants int
	Tunnels [][]int,
	Occupied bool
	
}

func main() {
	ants := 1

	v := make(map[string][][]int)

	v["A"] = [][]int{
		{0, 0, 3},
		{2, 2, 5},
		{3, 4, 0},
		{1, 8, 3},
	}



}
