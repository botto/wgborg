package main

import (
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
)

import (
	//#include <unistd.h>
	//#include <errno.h>
	"C"
)

// WGMgr container struct
type WGMgr struct {
	wgInt             *netlink.GenericLink
	wgClient          *wgctrl.Client
	store             *Store
	closing           chan bool
	rpcClient         *rpc.Client
	serverMode        bool
	shutdownFunctions []func()
}

func main() {
	initConfig()
	closingChan := make(chan bool)
	handleSignals(closingChan)
	wgMgr := WGMgr{
		closing:           closingChan,
		shutdownFunctions: make([]func(), 0),
	}
	serverFlag := flag.Bool("server", false, "set to run rpc server, otherwise assume http server")
	flag.Parse()
	if *serverFlag {
		startRPCServer(&wgMgr)
	} else {
		startHTTPServer(&wgMgr)
	}
	for {
		select {
		case <-wgMgr.closing:
			// Call each function that has registered a "shutdown" hooko
			for _, cb := range wgMgr.shutdownFunctions {
				cb()
			}
			wgMgr.cleanUp()
			os.Exit(0)
		}
	}
}

func (w *WGMgr) cleanUp() {
	if w.serverMode {
		err := netlink.LinkDel(w.wgInt)
		if err != nil {
			log.Fatal(err.Error())
		}
	} else {
		w.store.Close()
	}
}

func (w *WGMgr) setupStore() {
	storeCfg := DBServerConfig{
		Host:     config.DBHost,
		Port:     config.DBPort,
		User:     config.DBUser,
		Password: config.DBPassword,
		DBName:   config.DBName,
	}
	w.store = NewStore()
	w.store.Connect(&storeCfg)
}

// AddShutdownCB allow graceful shutdown
func (w *WGMgr) AddShutdownCB(cb func()) {
	w.shutdownFunctions = append(w.shutdownFunctions, cb)
}

// Handle termination of the application. Perform any cleanup required here.
func handleSignals(closing chan<- bool) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Sigterm caught, closing cleanely")
		closing <- true
	}()
}
