package main

func IsOccupied(roompos []int, antpos [][]int) bool {
	for _, ap := range antpos {
		if ap[0] == roompos[0] && ap[1] == roompos[1] {
			return true
		}
	}
	return false
}

func Canmove(room *Room, antpos [][]int) bool {
	if room.Occupied == false {
		room.Tunnels = nil
	}
	return true
}
