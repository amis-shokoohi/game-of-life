wasm:
	GOOS=js GOARCH=wasm go build -o http/main.wasm wasm/main.go
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" http

tinywasm:
	tinygo build -o http/main.wasm -target wasm --no-debug wasm/main.go
	cp "$(shell tinygo env GOROOT)/targets/wasm_exec.js" http

serve:
	go run dev-server/main.go