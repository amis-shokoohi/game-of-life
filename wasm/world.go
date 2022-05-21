package main

import (
	"math/rand"
	"syscall/js"
	"time"
)

// World is a 2D grid
type World struct {
	n          int
	resolution int
	gen        [][]cell
	ctx2d      js.Value
}

type cell struct {
	currentState bool
	futureState  bool
}

// NewWorld creates a new world and initializes first generation
func NewWorld(n int, resolution int, ctx2d js.Value) *World {
	rand.Seed(time.Now().UnixNano())
	n = n / resolution
	// Create first generation
	gen := make([][]cell, n, n)
	for z := 0; z < n*n; z++ {
		y := z / n
		x := z % n
		if x == 0 {
			gen[y] = make([]cell, n, n)
		}
		gen[y][x] = createCell()
	}
	return &World{
		n:          n,
		resolution: resolution,
		gen:        gen,
		ctx2d:      ctx2d,
	}
}

func createCell() cell {
	if rand.Intn(10)+1 < 2 { // chance of creation = 20%
		return cell{currentState: true}
	}
	return cell{currentState: false}
}

// Evolve creates next generation
func (w *World) Evolve() {
	for z := 0; z < w.n*w.n; z++ {
		y := z / w.n
		x := z % w.n
		top := (y - 1 + w.n) % w.n
		currY := (y + w.n) % w.n
		bottom := (y + 1 + w.n) % w.n
		left := (x - 1 + w.n) % w.n
		currX := (x + w.n) % w.n
		right := (x + 1 + w.n) % w.n
		// Find number of neighbors
		neighbors := 0
		if w.gen[top][left].currentState {
			neighbors++
		}
		if w.gen[top][currX].currentState {
			neighbors++
		}
		if w.gen[top][right].currentState {
			neighbors++
		}
		if w.gen[currY][left].currentState {
			neighbors++
		}
		if w.gen[currY][right].currentState {
			neighbors++
		}
		if w.gen[bottom][left].currentState {
			neighbors++
		}
		if w.gen[bottom][currX].currentState {
			neighbors++
		}
		if w.gen[bottom][right].currentState {
			neighbors++
		}
		// GoL rules
		if neighbors == 3 {
			w.gen[y][x].futureState = true
		} else if neighbors == 2 && w.gen[y][x].currentState {
			w.gen[y][x].futureState = true
		} else {
			w.gen[y][x].futureState = false
		}
	}
}

// Paint updates canvas & copy nextGen to currGen
func (w *World) Paint() {
	for z := 0; z < w.n*w.n; z++ {
		y := z / w.n
		x := z % w.n
		// Paint current cell
		if w.gen[y][x].currentState {
			w.ctx2d.Call("fillRect", x*w.resolution, y*w.resolution, w.resolution, w.resolution)
		} else {
			w.ctx2d.Call("clearRect", x*w.resolution, y*w.resolution, w.resolution, w.resolution)
		}
		// Copy new cell to current cell
		w.gen[y][x].currentState = w.gen[y][x].futureState
	}
}
