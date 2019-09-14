package main

import (
	"net/http"
)

func handlerAddPeer(w http.ResponseWriter, r *http.Request) {
	// var newPeerData Peer
	// err := json.NewDecoder(r.Body).Decode(&newPeerData)

	// if err != nil {
	// 	http.Error(w, err.Error(), 400)
	// 	return
	// }
	// if len(newPeerData.PublicKey) != 44 {
	// 	http.Error(w, "Public key is not 44 char long", 400)
	// 	return
	// }
	// if len(newPeerData.Psk) != 44 {
	// 	http.Error(w, "PSK is not 44 cahr long", 400)
	// 	return
	// }
	// if len(newPeerData.Name) == 0 || len(newPeerData.Name) > 255 {
	// 	http.Error(w, "Name must be > 0 or < 255", 400)
	// 	return
	// }
	// _, _, err = net.ParseCIDR(newPeerData.IP)
	// if len(newPeerData.IP) < 10 || len(newPeerData.IP) > 20 || err != nil {
	// 	http.Error(w, "IP Address could not be parsed as CIDR address (i.e.: 123.123.123.123/128)", 400)
	// 	return
	// }
	// newPeerSQL := `
	// 	INSERT INTO peers (peer_name, public_key, psk, ip)
	// 	VALUES ($1, $2, $3, $4)`
	// var insertErr error
	// _, insertErr = db.Exec(newPeerSQL, newPeerData.Name, newPeerData.PublicKey, newPeerData.Psk, newPeerData.IP)
	// if insertErr != nil {
	// 	http.Error(w, insertErr.Error(), 400)
	// 	log.Fatal(insertErr)
	// }
	// loadPeers()
}
