package main

import (
	"context"
	"log"
	"net/http"
	"net/rpc"
)

func startHTTPServer(wgMgr *WGMgr) {
	srv := &http.Server{Addr: ":8080"}
	var err error
	//make connection to rpc server
	wgMgr.rpcClient, err = rpc.DialHTTP("tcp", "127.0.0.1:39252")
	if err != nil {
		log.Fatalf("Error in dialing. %s", err)
	}
	wgMgr.setupStore()
	wgMgr.setupInterfaces()
	wgMgr.setupRoutes()
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()
	wgMgr.AddShutdownCB(func() {
		srv.Shutdown(context.Background())
	})
}

func (w *WGMgr) setupRoutes() {
	http.HandleFunc("/add_peer", w.handlerAddPeer)
}
