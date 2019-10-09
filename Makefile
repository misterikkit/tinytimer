
wasm: wasm_site/main.wasm
	date > wasm_site/built_at.txt

wasm_site/main.wasm: main.go
	GOOS=js GOARCH=wasm go build -o wasm_site/main.wasm main.go

clean:
	rm wasm_site/main.wasm