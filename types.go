package main

import (
	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Peer is a WireGuard peer
type Peer struct {
	ID        string `json:"id,omitempty" validate:"omitempty,uuid"`
	PublicKey string `json:"public_key" validate:"required,len=44"`
	Name      string `json:"name" validate:"required,gte=1,lte=255"`
	Psk       string `json:"psk" validate:"required,len=44"`
	IP        string `json:"ip" validate:"required,cidr"`
	NetworkID string `json:"network" validate:"required,uuid"`
}

// Network contains peers
type Network struct {
	ID         string `json:"id,omitempty" validate:"omitempty,uuid"`
	Name       string `json:"name" validate:"required,gte=1,lte=255"`
	PrivateKey string `json:"private_key" validate:"required,len=44"`
	Port       int    `json:"port" validte:"gte=1024,lte=65535"`
	IP         string `json:"ip" validate:"required,cidr"`
}

// InterfacePeersConfig is the peers list of the interface
type InterfacePeersConfig struct {
	WGPeers       *[]wgtypes.PeerConfig
	InterfaceName string
}

// WGInterface Internal representation of WG interface
type WGInterface struct {
	ID        string
	Interface *netlink.GenericLink
}
