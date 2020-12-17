package main

import (
	"math/rand"
	"syscall/js"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	document := js.Global().Get("document")

	width := js.Global().Get("innerWidth").Int() - 2
	height := js.Global().Get("innerHeight").Int() - 2

	length := height
	if height > width {
		length = width
	}

	res := 20
	length = length - length%res

	cnvs := document.Call("createElement", "canvas")
	cnvs.Set("width", length)
	cnvs.Set("height", length)
	cnvs.Get("style").Set("border", "1px solid #464658")
	document.Get("body").Call("appendChild", cnvs)

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

	var tMaxFPS float64 = 1000 / 60
	var repaint js.Func
	var lastTimestamp float64
	repaint = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		timestamp := args[0].Float()
		if timestamp-lastTimestamp >= tMaxFPS {
			evolve(&currGen, &nextGen, length)
			for y := 0; y < length; y++ {
				for x := 0; x < length; x++ {
					// Paint current cell
					if currGen[y][x] == 1 {
						ctx2d.Call("fillRect", x*res, y*res, res, res)
					} else {
						ctx2d.Call("clearRect", x*res, y*res, res, res)
					}
					// Copy new cell to current cell
					currGen[y][x] = nextGen[y][x]
				}
			}
			lastTimestamp = timestamp
		}
		js.Global().Call("requestAnimationFrame", repaint)
		return nil
	})
	defer repaint.Release()
	js.Global().Call("requestAnimationFrame", repaint)

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
