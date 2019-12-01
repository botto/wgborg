.PHONY: build debug
all: build

build:
	go build -ldflags=-w -o build/wgmgr

build_debug:
	go build -gcflags=all="-N -l" -o build/wgmgr

run: build
	./build/wgmgr

run_server: build
	sudo ./build/wgmgr -server

debug: build_debug
	sudo gdb ./build/wgmgr

debug_server: build_debug
	sudo gdb --args ./build/wgmgr -server

