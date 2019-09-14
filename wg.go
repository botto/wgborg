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

// SetupInterfaces gets all interface from the store
// and configurs them and adds the peers.
func (w *WGMgr) SetupInterfaces() {
	var err error
	wgClient, err = wgctrl.New()
	if err != nil {
		log.Fatalf("failed to open wgctrl: %v", err)
	}
	networks, err := w.store.LoadNetworks()
	if err != nil {
		fmt.Printf("Err :%s", err)
		log.Printf("Could not load networks.")
		log.Print(err)
		return
	}
	for _, ifDev := range networks {
		configureInterface(ifDev.Port, ifDev.PrivateKey, ifDev.Name)
		interfacePeers := w.GetNetworkPeers(ifDev.ID)
		addWgPeersToDevice(interfacePeers, ifDev.Name)
	}
}

func configureInterface(port int, priveKeyRaw string, name string) {
	privKey, err := wgtypes.ParseKey(priveKeyRaw)
	if err != nil {
		log.Printf("Could not parse key ")
	}
	// Generate new netlink of wireguard type
	linkAttrs := netlink.NewLinkAttrs()
	linkAttrs.Name = name
	wgInt = &netlink.GenericLink{LinkAttrs: linkAttrs, LinkType: "wireguard"}
	err = netlink.LinkAdd(wgInt)
	if err != nil {
		log.Fatalf("could not add '%s' (%v)\n", linkAttrs.Name, err)
	}
	deviceConfig := wgtypes.Config{
		PrivateKey: &privKey,
		ListenPort: &port,
	}
	err = wgClient.ConfigureDevice(wgInt.Name, deviceConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func peerToWgPeer(peerDef Peer) (*wgtypes.PeerConfig, error) {
	_, ipv4Net, err := net.ParseCIDR(peerDef.IP)
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

func addWgPeersToDevice(wgPeers *[]wgtypes.PeerConfig, deviceName string) {
	deviceConfig := wgtypes.Config{
		ReplacePeers: true,
		Peers:        *wgPeers,
	}
	wgClient.ConfigureDevice(deviceName, deviceConfig)
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
