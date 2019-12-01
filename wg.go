package main

import (
	"fmt"
	"log"
	"net"

	uuid "github.com/google/uuid"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (w *WGMgr) setupClient() {
	var err error
	w.wgClient, err = wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
}

// SetupInterfaces gets all interface from the store
// and configurs them and adds the peers.
func (w *WGMgr) setupInterfaces() {
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
		}
		var rpcRes interface{}
		fmt.Println("Calling RPC")
		w.rpcClient.Call("WGRpc.ConfigureInterface", iConfig, rpcRes)
		interfacePeers := w.GetNetworkPeers(ifDev.ID)
		peersConfig := &InterfacePeersConfig{
			WGPeers:       interfacePeers,
			InterfaceName: ifDev.Name,
		}
		w.rpcClient.Call("WGRpc.SetPeersOnInterface", peersConfig, rpcRes)
	}
}

// ConfigureInterface sets up the wg interface
func (w *WGMgr) ConfigureInterface(iConfig *InterfaceConfig) error {
	privKey, err := wgtypes.ParseKey(iConfig.PrivateKeyString)
	if err != nil {
		return fmt.Errorf("Could not parse key %s", err)
	}
	// Generate new netlink of wireguard type
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = iConfig.InterfaceName
	w.wgInt = &netlink.GenericLink{LinkAttrs: linkAttrs, LinkType: "wireguard"}
	err = netlink.LinkAdd(w.wgInt)
	if err != nil {
		return fmt.Errorf("Could not add '%s' (%v)", linkAttrs.Name, err)
	}
	deviceConfig := wgtypes.Config{
		PrivateKey: &privKey,
		ListenPort: &iConfig.Port,
	}
	err = w.wgClient.ConfigureDevice(linkAttrs.Name, deviceConfig)
	if err != nil {
		return fmt.Errorf("[WG/ConfigureInterface] Could not configure device '%s' (%v)", linkAttrs.Name, err)
	}
	return nil
}

// SetPeersOnInterface sets interface to only have the set peers
func (w *WGMgr) SetPeersOnInterface(conf *InterfacePeersConfig) error {
	deviceConfig := wgtypes.Config{
		ReplacePeers: true,
		Peers:        *conf.WGPeers,
	}
	err := w.wgClient.ConfigureDevice(conf.InterfaceName, deviceConfig)
	if err != nil {
		return fmt.Errorf("[WG/SetPeersOnInterface] Could not set peers on device '%s' (%v)", conf.InterfaceName, err)
	}
	return nil
}

// AddWgPeersToInterface adds peer to existing list of devices
func (w *WGMgr) AddWgPeersToInterface(conf *InterfacePeersConfig) error {
	deviceConfig := wgtypes.Config{
		ReplacePeers: false,
		Peers:        *conf.WGPeers,
	}
	err := w.wgClient.ConfigureDevice(conf.InterfaceName, deviceConfig)
	if err != nil {
		return fmt.Errorf("[WG/AddWgPeersToInterface] Could not add peers on device '%s' (%v)", conf.InterfaceName, err)
	}
	return nil
}

// GetNetworkPeers grabs peers from DB.
func (w *WGMgr) GetNetworkPeers(networkID *uuid.UUID) *[]wgtypes.PeerConfig {
	var wgPeers []wgtypes.PeerConfig
	peers, err := w.store.LoadPeers(networkID.String())
	fmt.Printf("Peers: %s\n", peers)
	if err != nil {
		log.Fatalf("Could not load peers %s", err)
	}
	for _, p := range peers {
		wgP, err := peerToWgPeer(p)
		if err != nil {
			log.Fatalf("Could not define WG peer, %s", err)
		}
		wgPeers = append(wgPeers, *wgP)
	}
	return &wgPeers
}

func peerToWgPeer(peerDef Peer) (*wgtypes.PeerConfig, error) {
	_, ipv4Net, err := net.ParseCIDR(peerDef.CIDR)
	if err != nil {
		log.Printf("Could not parse cidr for peer %s", err)
		return nil, err
	}
	pubKey, err := wgtypes.ParseKey(peerDef.PublicKey)
	if err != nil {
		log.Printf("Could not parse public key for peer %s", err)
		return nil, err
	}
	psk, err := wgtypes.ParseKey(peerDef.Psk)
	if err != nil {
		log.Printf("Could not parse public key for peer %s", err)
		return nil, err
	}
	wgPeer := wgtypes.PeerConfig{
		PublicKey:    pubKey,
		PresharedKey: &psk,
		AllowedIPs:   []net.IPNet{*ipv4Net},
	}
	return &wgPeer, nil
}