package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DBServerConfig database connection config
type DBServerConfig struct {
	Host     string
	Port     uint32
	User     string
	Password string
	DBName   string
}

// Store for network and peer config
type Store struct {
	db      *sql.DB
	closing chan bool
}

// NewStore set up new store
func NewStore() *Store {
	db := Store{
		closing: make(chan bool),
	}
	return &db
}

// Close does steps to close Store connection
func (s *Store) Close() {
	s.db.Close()
}

// Connect to the DB and check it works
func (s *Store) Connect(cfg *DBServerConfig) {
	plsqlInfo := fmt.Sprintf(`host=%s port=%d
		user=%s password=%s dbname=%s sslmode=disable`,
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	var err error
	s.db, err = sql.Open("postgres", plsqlInfo)
	if err != nil {
		log.Fatalf("Could not connect to DB, is it running? \nErr: %s", err)
	}

	err = s.db.Ping()
	if err != nil {
		log.Fatalf("Could connect to DB but not ping, strange. \nErr: %s", err)
	}
}

// LoadPeers gets peers from the store
func (s *Store) LoadPeers(networkID string) ([]Peer, error) {
	var newPeers []Peer
	allPeersSQL := `
		SELECT
			peer_name,
			public_key,
			psk,
			cidr
		FROM peers
		WHERE network = $1`
	rows, err := s.db.Query(allPeersSQL, networkID)
	if err != nil {
		log.Printf("Could not query. \n Err: %s", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		newPeer := Peer{}
		err = rows.Scan(
			&newPeer.Name,
			&newPeer.PublicKey,
			&newPeer.Psk,
			&newPeer.IP,
		)
		newPeers = append(newPeers, newPeer)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if len(newPeers) > 0 {
		return newPeers, nil
	}
	return nil, errors.New("no peers found")
}

// LoadNetworks gets the networks from the store and returns as Network list.
func (s *Store) LoadNetworks() ([]Network, error) {
	newNetworks := []Network{}
	allPeersSQL := `
		SELECT
			id,
			network_name,
			private_key,
			port,
			cidr
		FROM networks`
	rows, err := s.db.Query(allPeersSQL)
	if err != nil {
		log.Print("Could not query")
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		newNetwork := Network{}
		err = rows.Scan(
			&newNetwork.ID,
			&newNetwork.Name,
			&newNetwork.PrivateKey,
			&newNetwork.Port,
			&newNetwork.CIDR,
		)
		newNetworks = append(newNetworks, newNetwork)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	if len(newNetworks) > 0 {
		return newNetworks, nil
	}
	return nil, errors.New("No networks found")
}

// AddPeer add a new peer to a specific network
func (s *Store) AddPeer(newPeer Peer) {
	newPeerSQL := `
		INSERT INTO peers (peer_name, public_key, psk, ip, network)
		VALUES ($1, $2, $3, $4)`
	_, err := s.db.Exec(
		newPeerSQL,
		newPeer.Name,
		newPeer.PublicKey,
		newPeer.Psk,
		newPeer.IP,
		newPeer.NetworkID,
	)
	if err != nil {
		log.Fatal(err)
	}
}

// func (s *Store) AddNetwork(newNetwork Network) {
// 	newNetworkSQL := `
// 		INSERT INTO networks (name, private_key, cidr)
// 		VALUES ($1, $2, $3)`
// 	_, err := db.Exec(newNetworkSQL, newNetwork.Name, newNetwork.PrivateKey, newNetwork.CIDR)
// 	if err != nil {
// 		http.Error(w, insertErr.Error(), 400)
// 		log.Fatal(insertErr)
// 	}
// }
