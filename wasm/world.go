package main

import (
	"math/rand"
	"syscall/js"
	"time"
)

// World is a 2D world
type World struct {
	length  int
	currGen *[][]uint8
	nextGen *[][]uint8
	ctx2d   *js.Value
}

// NewWorld creates a new world and initializes first generation
func NewWorld(length int, ctx2d *js.Value) *World {
	rand.Seed(time.Now().UnixNano())
	// Allocate two 2D arrays
	currGen := make([][]uint8, length, length)
	nextGen := make([][]uint8, length, length)
	// Create first generation
	for y := 0; y < length; y++ {
		currGen[y] = make([]uint8, length, length)
		nextGen[y] = make([]uint8, length, length)
		for x := 0; x < length; x++ {
			currGen[y][x] = createCell()
		}
	}
	return &World{
		length:  length,
		currGen: &currGen,
		nextGen: &nextGen,
		ctx2d:   ctx2d,
	}
}

func createCell() uint8 {
	chance := 20 // percent
	if rand.Intn(100)+1 < chance {
		return 1
	}
	return 0
}

// Evolve creates next generation
func (w *World) Evolve() {
	for y := 0; y < w.length; y++ {
		top := (y - 1 + w.length) % w.length
		currY := (y + w.length) % w.length
		bottom := (y + 1 + w.length) % w.length
		for x := 0; x < w.length; x++ {
			left := (x - 1 + w.length) % w.length
			currX := (x + w.length) % w.length
			right := (x + 1 + w.length) % w.length
			neighbors := 0
			// Find number of neighbors
			if (*w.currGen)[top][left] == 1 {
				neighbors++
			}
			if (*w.currGen)[top][currX] == 1 {
				neighbors++
			}
			if (*w.currGen)[top][right] == 1 {
				neighbors++
			}
			if (*w.currGen)[currY][left] == 1 {
				neighbors++
			}
			if (*w.currGen)[currY][right] == 1 {
				neighbors++
			}
			if (*w.currGen)[bottom][left] == 1 {
				neighbors++
			}
			if (*w.currGen)[bottom][currX] == 1 {
				neighbors++
			}
			if (*w.currGen)[bottom][right] == 1 {
				neighbors++
			}
			// GoL rules
			if neighbors == 3 {
				(*w.nextGen)[y][x] = 1
			} else if neighbors == 2 && (*w.currGen)[y][x] == 1 {
				(*w.nextGen)[y][x] = 1
			} else {
				(*w.nextGen)[y][x] = 0
			}
		}
	}
}

// Paint updates canvas & copy nextGen to currGen
func (w *World) Paint(res int) {
	for y := 0; y < w.length; y++ {
		for x := 0; x < w.length; x++ {
			// Paint current cell
			if (*w.currGen)[y][x] == 1 {
				w.ctx2d.Call("fillRect", x*res, y*res, res, res)
			} else {
				w.ctx2d.Call("clearRect", x*res, y*res, res, res)
			}
			// Copy new cell to current cell
			(*w.currGen)[y][x] = (*w.nextGen)[y][x]
		}
	}
}
