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

func (w *WGMgr) setupClient() {
	var err error
	w.wgClient, err = wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
}

// initInterfaces gets all interface from the store
// and configurs them and adds the peers.
func (w *WGMgr) initInterfaces() {
	networks, err := w.store.LoadNetworks()
	if err != nil {
		fmt.Printf("Err :%s", err)
		log.Printf("Could not load networks.")
		log.Print(err)
		return
	}
	for _, ifDev := range networks {
		iConfig := &InterfaceConfig{
			Port:             ifDev.Port,
			PrivateKeyString: ifDev.PrivateKey,
			InterfaceName:    ifDev.Name,
			IP:               ifDev.IP,
		}
		var rpcRes interface{}
		w.rpcClient.Call("WGRpc.ConfigureInterface", iConfig, rpcRes)
		interfacePeers := w.GetNetworkPeers(ifDev.ID)
		peersConfig := InterfacePeersConfig{
			WGPeers:       interfacePeers,
			InterfaceName: ifDev.Name,
		}
		w.rpcClient.Call("WGRpc.SetPeersOnInterface", &peersConfig, rpcRes)
	}
}
