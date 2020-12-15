package main

import (
	"math/rand"
	"syscall/js"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	document := js.Global().Get("document")

	width := js.Global().Get("innerWidth").Int()
	height := js.Global().Get("innerHeight").Int()
	var length int
	if height > width {
		length = width
	} else {
		length = height
	}

	res := 20
	length = length - length%res

	container := document.Call("querySelector", ".container")
	container.Get("style").Set("width", length)
	container.Get("style").Set("height", length)
	container.Get("style").Set("border", "1px solid #464658")

	cnvs := document.Call("getElementById", "cnvs")
	cnvs.Set("width", length)
	cnvs.Set("height", length)
	ctx2d := cnvs.Call("getContext", "2d")
	ctx2d.Set("fillStyle", "#78ABB7")

	length = length / res

	// Allocate two 2D arrays
	currGen := make([][]uint8, length, length)
	nextGen := make([][]uint8, length, length)
	for i := 0; i < length; i++ {
		currGen[i] = make([]uint8, length, length)
		nextGen[i] = make([]uint8, length, length)
	}

	// Create first generation
	for y := 0; y < length; y++ {
		for x := 0; x < length; x++ {
			currGen[y][x] = createCell()
		}
	}

	select {}
}

func createCell() uint8 {
	chance := 20 // percent
	if rand.Intn(100)+1 < chance {
		return 1
	}
	return 0
}

func evolve(currGen *[][]uint8, nextGen *[][]uint8, length int) {
	for y := 0; y < length; y++ {
		for x := 0; x < length; x++ {
			neighbors := 0
			// Find number of neighbors
			if (*currGen)[(y-1+length)%length][(x-1+length)%length] == 1 { // top left
				neighbors++
			}
			if (*currGen)[(y-1+length)%length][(x+length)%length] == 1 { // top
				neighbors++
			}
			if (*currGen)[(y-1+length)%length][(x+1+length)%length] == 1 { // top right
				neighbors++
			}
			if (*currGen)[(y+length)%length][(x-1+length)%length] == 1 { // left
				neighbors++
			}
			if (*currGen)[(y+length)%length][(x+1+length)%length] == 1 { // right
				neighbors++
			}
			if (*currGen)[(y+1+length)%length][(x-1+length)%length] == 1 { // bottom left
				neighbors++
			}
			if (*currGen)[(y+1+length)%length][(x+length)%length] == 1 { // bottom
				neighbors++
			}
			if (*currGen)[(y+1+length)%length][(x+1+length)%length] == 1 { // bottom right
				neighbors++
			}
			// GoL rules
			if neighbors == 3 {
				(*nextGen)[y][x] = 1
			} else if neighbors == 2 && (*currGen)[y][x] == 1 {
				(*nextGen)[y][x] = 1
			} else {
				(*nextGen)[y][x] = 0
			}
		}
	}
}
