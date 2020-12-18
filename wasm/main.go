package main

import (
	"syscall/js"
)

func main() {
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

	world := NewWorld(length, &ctx2d)

	var tMaxFPS float64 = 1000 / 60
	var repaint js.Func
	var lastTimestamp float64
	repaint = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		timestamp := args[0].Float()
		if timestamp-lastTimestamp >= tMaxFPS {
			world.Evolve()
			world.Paint(res)
			lastTimestamp = timestamp
		}
		js.Global().Call("requestAnimationFrame", repaint)
		return nil
	})
	defer repaint.Release()
	js.Global().Call("requestAnimationFrame", repaint)

	select {}
}
