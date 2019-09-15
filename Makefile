.PHONY: build debug
all: build

build:
	go build -ldflags=-w -o build/wgmgr 

debug:
	go build -gcflags=all="-N -l" -o build/wgmgr
	sudo gdb ./build/wgmgr