package main

import (
	"log"
)

// WGRpc is the RPC interface to expose.
type WGRpc struct {
	wgMgr *WGMgr
}

// ConfigReponseMsg is a simple int.
type ConfigReponseMsg int

// ConfigureInterface set up WG interface.
func (rpc *WGRpc) ConfigureInterface(iConfig InterfaceConfig, reply *int) error {
	err := rpc.wgMgr.ConfigureInterface(&iConfig)
	if err != nil {
		log.Fatalf("[RPC/ConfigureInterface] Err: %s", err)
		return err
	}
	*reply = 1
	return nil
}

// SetPeersOnInterface sets peers on a WG Interface.
func (rpc *WGRpc) SetPeersOnInterface(iConfig InterfacePeersConfig, reply *int) error {
	err := rpc.wgMgr.SetPeersOnInterface(&iConfig)
	if err != nil {
		log.Fatalf("[RPC/SetPeersOnInterface] Err: %s", err)
		return err
	}
	*reply = 1
	return nil
}

// AddWgPeersToInterface sets peers on a WG Interface.
func (rpc *WGRpc) AddWgPeersToInterface(iConfig InterfacePeersConfig, reply *int) error {
	err := rpc.wgMgr.AddWgPeersToInterface(&iConfig)
	if err != nil {
		log.Fatalf("[RPC/AddWgPeersToInterface] Err: %s", err)
		return err
	}
	*reply = 1
	return nil
}