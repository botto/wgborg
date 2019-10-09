.PHONY: build debug
all: build

build:
	go build -ldflags=-w -o build/wgmgr

build_debug:
	go build -gcflags=all="-N -l" -o build/wgmgr

run: build
	sudo ./build/wgmgr

debug: build_debug
	sudo gdb ./build/wgmgr

debug_server: build_debug
	sudo gdbserver 127.0.0.1:33333 ./build/wgmgr
