
wasm: _out/main.wasm

_out/main.wasm: main.go
	GOOS=js GOARCH=wasm go build -o _out/main.wasm main.go