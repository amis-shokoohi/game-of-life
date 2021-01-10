package main

import (
	"math/rand"
	"syscall/js"
	"time"
)

// World is a 2D world
type World struct {
	length int
	gen    [][]cell
	ctx2d  *js.Value
}

type cell struct {
	currentState bool
	futureState  bool
}

// NewWorld creates a new world and initializes first generation
func NewWorld(length int, ctx2d *js.Value) *World {
	rand.Seed(time.Now().UnixNano())
	// Create first generation
	gen := make([][]cell, length, length)
	for y := 0; y < length; y++ {
		gen[y] = make([]cell, length, length)
		for x := 0; x < length; x++ {
			gen[y][x] = createCell()
		}
	}
	return &World{
		length: length,
		gen:    gen,
		ctx2d:  ctx2d,
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
	for y := 0; y < w.length; y++ {
		top := (y - 1 + w.length) % w.length
		currY := (y + w.length) % w.length
		bottom := (y + 1 + w.length) % w.length
		for x := 0; x < w.length; x++ {
			left := (x - 1 + w.length) % w.length
			currX := (x + w.length) % w.length
			right := (x + 1 + w.length) % w.length
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
func (w *World) Paint(res int) {
	for y := 0; y < w.length; y++ {
		for x := 0; x < w.length; x++ {
			// Paint current cell
			if w.gen[y][x].currentState {
				w.ctx2d.Call("fillRect", x*res, y*res, res, res)
			} else {
				w.ctx2d.Call("clearRect", x*res, y*res, res, res)
			}
			// Copy new cell to current cell
			w.gen[y][x].currentState = w.gen[y][x].futureState
		}
	}
}
