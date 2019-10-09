package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (wg *WGMgr) handlerAddPeer(w http.ResponseWriter, r *http.Request) {
	var newPeerData Peer
	err := json.NewDecoder(r.Body).Decode(&newPeerData)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if len(newPeerData.PublicKey) != 44 {
		http.Error(w, "Public key is not 44 char long", 400)
		return
	}
	if len(newPeerData.Psk) != 44 {
		http.Error(w, "PSK is not 44 cahr long", 400)
		return
	}
	if len(newPeerData.Name) == 0 || len(newPeerData.Name) > 255 {
		http.Error(w, "Name length must be > 0 or < 255", 400)
		return
	}
	if len(newPeerData.NetworkID) == 36 {
		http.Error(w, "Missing valid Network uuid.", 400)
		return
	}
	_, _, err = net.ParseCIDR(newPeerData.CIDR)
	if err != nil {
		http.Error(w, "IP Address could not be parsed as CIDR address (i.e.: 123.123.123.123/128)", 400)
		return
	}
	networkName, err := wg.store.GetNetworkNameByID(newPeerData.NetworkID)
	if err != nil {
		http.Error(w, "Network ID could not be found", 400)
		log.Printf("Could not find %s", err)
		return
	}
	wg.store.AddPeer(&newPeerData)
	newWgPeer, err := peerToWgPeer(newPeerData)
	if err != nil {
		http.Error(w, "Uho, something bad happend, we will look in to this", 400)
		return
	}
	newPeers := []wgtypes.PeerConfig{*newWgPeer}
	wg.addWgPeersToDevice(&newPeers, networkName)
}

func (wg *WGMgr) handlerAddNetwork(w http.ResponseWriter, r *http.Request) {
	var newNetworkData Network
	err := json.NewDecoder(r.Body).Decode(&newNetworkData)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if len(newNetworkData.PrivateKey) != 44 {
		http.Error(w, "Private key is not 44 char long", 400)
		return
	}
	if newNetworkData.Port < 1024 || newNetworkData.Port > 65535 {
		http.Error(w, "Port range is not within 1024 -> 65535", 400)
		return
	}
	if len(newNetworkData.Name) == 0 || len(newNetworkData.Name) > 255 {
		http.Error(w, "Name length must be > 0 or < 255", 400)
		return
	}
	_, _, err = net.ParseCIDR(newNetworkData.CIDR)
	if err != nil {
		http.Error(w, "IP Address could not be parsed as CIDR address (i.e.: 123.123.123.123/128)", 400)
		return
	}
}
