all: build
	./modbusHealthChecker.exe

build:
	go mod tidy
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	GOARCH=arm GOOS=linux go build -o modbusHealthChecker
	go build -o modbusHealthChecker.exe
clean:
	rm -f modbusHealthChecker.exe
	rm -r  modbusHealthChecker
	rm -f web/app.wasm
