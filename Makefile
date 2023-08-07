all: build
	./pilones.com.exe

build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build
clean:
	rm -f pilones.com.exe
	rm -f web/app.wasm
