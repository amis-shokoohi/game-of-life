package main

import (
	"math/rand"
	"syscall/js"
	"time"
)

// World is a 2D world
type World struct {
	width      int
	resolution int
	gen        [][]cell
	ctx2d      js.Value
}

type cell struct {
	currentState bool
	futureState  bool
}

// NewWorld creates a new world and initializes first generation
func NewWorld(width int, resolution int, ctx2d js.Value) *World {
	rand.Seed(time.Now().UnixNano())
	width = width / resolution
	// Create first generation
	gen := make([][]cell, width, width)
	for y := 0; y < width; y++ {
		gen[y] = make([]cell, width, width)
		for x := 0; x < width; x++ {
			gen[y][x] = createCell()
		}
	}
	return &World{
		width:      width,
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
	for y := 0; y < w.width; y++ {
		top := (y - 1 + w.width) % w.width
		currY := (y + w.width) % w.width
		bottom := (y + 1 + w.width) % w.width
		for x := 0; x < w.width; x++ {
			left := (x - 1 + w.width) % w.width
			currX := (x + w.width) % w.width
			right := (x + 1 + w.width) % w.width
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
}

// Paint updates canvas & copy nextGen to currGen
func (w *World) Paint() {
	for y := 0; y < w.width; y++ {
		for x := 0; x < w.width; x++ {
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
}
