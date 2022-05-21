package main

import (
	"syscall/js"
)

var FPS float64 = 60.0

func main() {
	document := js.Global().Get("document")

	width := js.Global().Get("innerWidth").Int() - 2
	height := js.Global().Get("innerHeight").Int() - 2
	if height < width {
		width = height
	}

	resolution := 20
	width = width - width%resolution

	canvas := document.Call("createElement", "canvas")
	canvas.Set("width", width)
	canvas.Set("height", width)
	canvas.Get("style").Set("border", "1px solid #464658")
	document.Get("body").Call("appendChild", canvas)

	ctx2d := canvas.Call("getContext", "2d")
	ctx2d.Set("fillStyle", "#78ABB7")

	world := NewWorld(width, resolution, ctx2d)

	var timeInterval float64 = 1000 / FPS
	var repaint js.Func
	var lastTimestamp float64
	repaint = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		timestamp := args[0].Float()
		if timestamp-lastTimestamp >= timeInterval {
			world.Evolve()
			world.Paint()
			lastTimestamp = timestamp
		}
		js.Global().Call("requestAnimationFrame", repaint)
		return nil
	})
	defer repaint.Release()
	js.Global().Call("requestAnimationFrame", repaint)

	select {}
}
