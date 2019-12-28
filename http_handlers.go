package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func (wg *WGMgr) handlerAddPeer(w http.ResponseWriter, r *http.Request) {
	var newPeerData Peer
	err := json.NewDecoder(r.Body).Decode(&newPeerData)
	if err != nil {
		http.Error(w, "fatal error adding peer", 500)
		log.Printf("Could not unmarshal JSON: %s", err)
		return
	}
	err = wg.validate.Struct(newPeerData)
	if err != nil {
		for _, sgl_err := range err.(validator.ValidationErrors) {
			errStr := fmt.Sprintf("error in payload, field: %s is not a valid %s", sgl_err.Field(), sgl_err.Tag())
			http.Error(w, errStr, 400)
		}
		return
	}
	networkName, err := wg.store.GetNetworkNameByID(newPeerData.NetworkID)
	if err != nil {
		errStr := fmt.Sprintf("network not found: %s", newPeerData.NetworkID)
		http.Error(w, errStr, 404)
		return
	}
	wg.store.AddPeer(&newPeerData)
	newWgPeer, err := peerToWgPeer(newPeerData)
	if err != nil {
		http.Error(w, "error adding peer", 500)
		log.Printf("Could not add peer: %s, err: %s", newPeerData, err.Error())
		return
	}
	newPeers := []wgtypes.PeerConfig{*newWgPeer}
	peersConfig := InterfacePeersConfig{
		WGPeers:       &newPeers,
		InterfaceName: networkName,
	}
	var rpcRes interface{}
	wg.rpcClient.Call("WGRpc.AddWgPeersToInterface", &peersConfig, rpcRes)
}

func (wg *WGMgr) handlerAddNetwork(w http.ResponseWriter, r *http.Request) {
	var newNetworkData Network
	var newID string
	err := json.NewDecoder(r.Body).Decode(&newNetworkData)
	if err != nil {
		log.Printf("Could not unmarshal JSON")
		http.Error(w, "fatal error adding peer", 500)
		return
	}
	err = wg.validate.Struct(newNetworkData)
	if err != nil {
		for _, sgl_err := range err.(validator.ValidationErrors) {
			errStr := fmt.Sprintf("error in payload, field: %s is not a valid %s", sgl_err.Field(), sgl_err.Tag())
			http.Error(w, errStr, 400)
		}
		return
	}

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	newID, err = wg.store.AddNetwork(&newNetworkData)
	if err != nil {
		http.Error(w, "There was an error adding the new network", 500)
		fmt.Printf("Error adding network: %s", err)
		return
	}
	newNetworkConfig := Network{
		ID:         newID,
		IP:         newNetworkData.IP,
		Port:       newNetworkData.Port,
		Name:       newNetworkData.Name,
		PrivateKey: newNetworkData.PrivateKey,
	}
	var rpcRes interface{}
	wg.rpcClient.Call("WGRpc.ConfigureInterface", newNetworkConfig, rpcRes)
	httpOut := map[string]interface{}{
		"NewID": newID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(httpOut)
}
