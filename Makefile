
wasm: wasm_site/main.wasm
	date > wasm_site/built_at.txt

wasm_site/main.wasm: main.go
	GOOS=js GOARCH=wasm go build -o wasm_site/main.wasm *.go

serve: wasm
	find -name '*.go' -or -name '*.js' -or -name '*.html' -or -name '*.css' | entr make clean wasm &
	goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`./wasm_site`)))'
	killall entr


clean:
	rm wasm_site/main.wasm