package main

func isOccupied(roomData [][]int, antpos [][]int) bool {
	if antpos == roomData["A"] || antpos == roomData["B"] || antpos == roomData["C"] || antpos == roomData["D"] {
		return true
	}
	return false
}
