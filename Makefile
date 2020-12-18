.PHONY: wasm tinywasm serve

wasm: wasm/main.go wasm/world.go
	GOOS=js GOARCH=wasm go build -o http/main.wasm ./wasm
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./http

tinywasm: wasm/main.go wasm/world.go
	tinygo build -o http/main.wasm -target wasm --no-debug ./wasm
	cp "$(shell tinygo env GOROOT)/targets/wasm_exec.js" ./http

serve: dev-server/main.go http/index.html
	go run ./dev-server/main.go