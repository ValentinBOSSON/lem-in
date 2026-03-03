package main

type chamber struct {
	ID          string
	coordinates []int
	occupied    bool
	tunnels     []*chamber
}
