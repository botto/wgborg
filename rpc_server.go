package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func startRPCServer(wgMgr *WGMgr) {
	wgMgr.serverMode = true
	wgMgr.setupClient()
	rpcServer := &WGRpc{
		wgMgr: wgMgr,
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
	wgMgr.AddShutdownCB(func() {
		l.Close()
	})
	go http.Serve(l, nil)
}
