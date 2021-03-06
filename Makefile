build: _out/flash.uf2

_out/flash.uf2: _out $(shell find -name '*.go')
	tinygo build -o _out/flash.uf2 -target=itsybitsy-m4

_out:
	mkdir _out

flash: build
	sudo mount -t drvfs e: /mnt/e
	cp _out/flash.uf2 -target=itsybitsy-m4

wasm: wasm_site/main.wasm
	date > wasm_site/built_at.txt
	tinygo version >> wasm_site/built_at.txt

wasm_site/main.wasm: $(shell find -name '*.go')
	tinygo build -o wasm_site/main.wasm -target=wasm

serve: wasm
	find -name '*.go' -or -name '*.js' -or -name '*.html' -or -name '*.css' | entr make clean wasm &
	goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`./wasm_site`)))'
	killall entr


clean:
	rm -f wasm_site/main.wasm
	rm -f _out/flash.uf2