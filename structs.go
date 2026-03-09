package main

type chamber struct {
	ID          string
	coordinates []int
	occupied    bool
	tunnels     []*chamber
}

type ant struct {
	ID           int
	coordination []int
	oldchambre   int
	currentRoom  *chamber
	assignedPath []*chamber
	pathPosition int
}

type Item struct {
	node     *chamber
	priority int
	index    int
}

type PriorityQueue []*Item

type RoomOccupancy struct {
	roomID    string
	timeFreed int
	antID     int
	index     int
}
