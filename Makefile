.PHONY: build debug
all: build

build:
	go build -o build/wg_mgr

debug: build
	sudo ./build/wg_mgr