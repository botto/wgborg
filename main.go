package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
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
	wgInt      *netlink.GenericLink
	wgClient   *wgctrl.Client
	store      *Store
	closing    chan bool
	rpcClient  *rpc.Client
	serverMode bool
}

func main() {
	initConfig()
	closingChan := make(chan bool)
	handleSignals(closingChan)
	wgMgr := WGMgr{
		closing: closingChan,
	}
	serverFlag := flag.Bool("server", false, "set to run rpc server, otherwise assume http server")
	flag.Parse()
	if *serverFlag {
		wgMgr.serverMode = true
		wgMgr.setupClient();
		rpcServer := &WGRpc{
			wgMgr: &wgMgr,
		}
		err := rpc.Register(rpcServer)
		if err != nil {
			log.Fatalf("Format of service WGMgr isn't correct. %s", err)
		}
		rpc.HandleHTTP()
		l, e := net.Listen("tcp", "127.0.0.1:39252")
		if e != nil {
			log.Fatalf("Couldn't start listening. Error %s", e)
		}
		log.Println("Serving RPC handler")
		go http.Serve(l, nil)
	} else {
		var err error
		//make connection to rpc server
		wgMgr.rpcClient, err = rpc.DialHTTP("tcp", "127.0.0.1:39252")
		if err != nil {
			log.Fatalf("Error in dialing. %s", err)
		}
		wgMgr.setupStore()
		wgMgr.setupInterfaces()
		wgMgr.setupRoutes()
		go setupServer()
	}
	for {
		select {
		case <-wgMgr.closing:
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

func (w *WGMgr) setupRoutes() {
	http.HandleFunc("/add_peer", w.handlerAddPeer)
}

func setupServer() {
	log.Fatal(http.ListenAndServe(":8080", nil))
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
