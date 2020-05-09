GOFMT=gofmt
GC=go build

VERSION := $(shell git describe --always --tags --long)
BUILD_NODE_PAR = -ldflags "-X main.Version=1.0.0"

ARCH=$(shell uname -m)
SRC_FILES = $(shell git ls-files | grep -e .go$ | grep -v _test.go)

ont-tool: $(SRC_FILES)
	$(GC)  $(BUILD_NODE_PAR) -o dht-tool main.go

dht-tool-cross: dht-tool-windows dht-tool-linux dht-tool-darwin

dht-tool-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o dht-tool-windows-amd64.exe main.go

dht-tool-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o dht-tool-linux-amd64 main.go

dht-tool-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GC) $(BUILD_NODE_PAR) -o dht-tool-darwin-amd64 main.go

tools-cross: tools-windows tools-linux tools-darwin

format:
	$(GOFMT) -w main.go

clean:
	rm -rf *.8 *.o *.out *.6 *exe
	rm -rf dht-tool dht-tool-*

restart:
	make clean && make dht-tool && ./dht-tool
