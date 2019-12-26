package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	uuid "github.com/google/uuid"
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
			ipv4
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
		return nil, err
	}
	if len(newPeers) > 0 {
		return newPeers, nil
	}
	return nil, nil
}

// LoadNetworks gets the networks from the store and returns as Network list.
func (s *Store) LoadNetworks() ([]Network, error) {
	newNetworks := []Network{}
	allPeersSQL := `
		SELECT
			id,
			name,
			private_key,
			port,
			ipv4
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
			&newNetwork.IP,
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
func (s *Store) AddPeer(newPeer *Peer) {
	newPeerSQL := `
		INSERT INTO peers (
			peer_name,
			public_key,
			psk,
			ipv4,
			network
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5
		)`
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

// AddNetwork adds new network
func (s *Store) AddNetwork(newNetwork *Network) (string, error) {
	var insertID string
	newNetworkPre, err := s.db.Prepare(`
		INSERT INTO networks (
			name,
			private_key,
			port,
			ipv4
		)
		VALUES (
			$1,
			$2,
			$3,
			$4
		)
		RETURNING id`)
	if err != nil {
		return "", fmt.Errorf("There was a problem preparing SQL: %s", err)
	}
	err = newNetworkPre.QueryRow(
		newNetwork.Name,
		newNetwork.PrivateKey,
		newNetwork.Port,
		newNetwork.IP,
	).Scan(&insertID)
	if err != nil {
		return "", fmt.Errorf("There was a problem executing the query: %s", err)
	}
	return insertID, nil
}

// GetNetworkNameByID returns the name of the WG network.
func (s *Store) GetNetworkNameByID(networkID *uuid.UUID) (string, error) {
	var name string
	networkNameSQL := `SELECT name FROM networks WHERE id=$1`
	row := s.db.QueryRow(networkNameSQL, networkID.String())
	err := row.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
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
