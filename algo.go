package main

func IsOccupied(x int, y int, antpos *ant) bool {
	if antpos.coordination[0] == x && antpos.coordination[1] == y {
		return true
	}
	return false
}

func Canmove(room *chamber, antpos [][]int) bool {
	if room.occupied == false {
		room.tunnels = nil
	}
	return true
}

/*func OptimizedPath(end *chamber) {
	if chamber.occupied == true {

	}
}
*/
