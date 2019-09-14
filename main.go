package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
)

// WGCtrl Client ref
var wgClient *wgctrl.Client

// WGInterface ref
var wgInt *netlink.GenericLink

// WGMgr container struct
type WGMgr struct {
	store   *Store
	closing chan bool
}

func main() {
	initConfig()
	closingChan := make(chan bool)
	handleSignals(closingChan)
	wgMgr := WGMgr{
		closing: closingChan,
	}
	wgMgr.setupStore()
	wgMgr.SetupInterfaces()
	//setupRoutes()
	for {
		select {
		case <-wgMgr.closing:
			wgMgr.cleanUp()
		}
	}
	// setupServer()
	// defer db.Close()
	// defer wgClient.Close()
	// defer cleanUp()
}

func (w *WGMgr) cleanUp() {
	w.store.Close()
	err := netlink.LinkDel(wgInt)
	fmt.Printf("Link %s\n", wgInt)
	if err != nil {
		log.Fatal(err.Error())
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

func setupRoutes() {
	http.HandleFunc("/add_peer", handlerAddPeer)
}

func setupServer() {
	//var err error
	// var devices []*wgtypes.Device
	// devices, err = wgClient.Devices()
	// if err != nil {
	// 	log.Fatalf("failed to get devices: %v", err)
	// }
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
		os.Exit(0)
	}()
}
