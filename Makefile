all: build wasm


flash.uf2: *.go
	tinygo build -o flash.uf2 -target=itsybitsy-m4 .

build: flash.uf2

mount:
	sudo mount -t drvfs e: /mnt/e

flash: mount flash.uf2
	cp flash.uf2 /mnt/e/

wasm: wasm_site/main.wasm
	date > wasm_site/built_at.txt

wasm_site/main.wasm: main.go
	GOOS=js GOARCH=wasm go build -tags=wasmsite -o wasm_site/main.wasm *.go

serve: wasm
	find -name '*.go' -or -name '*.js' -or -name '*.html' -or -name '*.css' | entr make clean wasm &
	goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`./wasm_site`)))'
	killall entr

clean:
	rm -f wasm_site/main.wasm