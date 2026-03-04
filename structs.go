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
}
