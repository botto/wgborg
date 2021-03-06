package main

import (
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// ConfigureInterface sets up the wg interface
func (w *WGMgr) ConfigureInterface(iConfig *Network) error {
	privKey, err := wgtypes.ParseKey(iConfig.PrivateKey)
	if err != nil {
		return fmt.Errorf("Could not parse key %s", err)
	}
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = iConfig.Name
	_, err = uuid.Parse(iConfig.ID)
	if err != nil {
		return fmt.Errorf("A valid UUID was not provided: %s", err)
	}
	newInterface := WGInterface{
		ID: iConfig.ID,
	}

	// Generate new interface through netlink (wireguard type)
	newInterface.Interface = &netlink.GenericLink{
		LinkAttrs: linkAttrs,
		LinkType:  "wireguard",
	}
	w.wgInt = append(w.wgInt, &newInterface)
	err = netlink.LinkAdd(newInterface.Interface)
	if err != nil {
		return fmt.Errorf("Could not add '%s' (%v)", linkAttrs.Name, err)
	}
	ipv4Addr, err := netlink.ParseAddr(iConfig.IP)
	if err != nil {
		return fmt.Errorf("IP is not valid IPv4 address %s", err)
	}
	err = netlink.AddrAdd(newInterface.Interface, ipv4Addr)
	if err != nil {
		return fmt.Errorf("Could not add addr %s to interface %s", iConfig.IP, err)
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
func (w *WGMgr) GetNetworkPeers(networkID string) *[]wgtypes.PeerConfig {
	var wgPeers []wgtypes.PeerConfig
	peers, err := w.store.LoadPeers(networkID)
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
	_, ipv4Net, err := net.ParseCIDR(peerDef.IP)
	if err != nil {
		log.Printf("IP is not valid cidr for peer %s", err)
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
