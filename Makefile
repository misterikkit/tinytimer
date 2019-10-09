
wasm: wasm_site/main.wasm
	date > wasm_site/built_at.txt

wasm_site/main.wasm: main.go
	GOOS=js GOARCH=wasm go build -o wasm_site/main.wasm main.go

serve: wasm
	goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`./wasm_site`)))'


clean:
	rm wasm_site/main.wasm