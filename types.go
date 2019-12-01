package main

import (
	uuid "github.com/google/uuid"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// Peer is a WireGuard peer
type Peer struct {
	PublicKey string     `json:"public_key"`
	Name      string     `json:"name"`
	Psk       string     `json:"psk"`
	CIDR      string     `json:"ip"`
	NetworkID *uuid.UUID `json:"network"`
}

// Network contains peers
type Network struct {
	ID         *uuid.UUID `json:"id,omitempty"`
	Name       string     `json:"name"`
	PrivateKey string     `json:"private_key"`
	Port       int        `json:"port"`
	CIDR       string     `json:"cidr"`
}

// InterfaceConfig is the interface configuration
type InterfaceConfig struct {
	Port             int
	PrivateKeyString string
	InterfaceName    string
}

// InterfacePeersConfig is the peers list of the interface
type InterfacePeersConfig struct {
	WGPeers       *[]wgtypes.PeerConfig
	InterfaceName string
}
