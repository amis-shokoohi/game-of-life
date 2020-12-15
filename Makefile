wasm: http
	GOOS=js GOARCH=wasm go build -o http/main.wasm main.go
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" http
	cp index.html http

tinywasm: http
	tinygo build -o http/main.wasm -target wasm --no-debug main.go
	cp "$(shell tinygo env GOROOT)/targets/wasm_exec.js" http
	cp index.html http

http:
	mkdir http