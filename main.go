package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"time"
)

type Client struct {
	PrivateKey    *ecdsa.PrivateKey
	PublicKeyHash string
	StartTime     time.Time
	EndTime       time.Time
	Contact       Contact
}

type Contact struct {
	PublicKey     *ecdsa.PublicKey
	PublicKeyHash string
}

func (c *Client) Start() {
	go func() {
		// get websocket

		// authenticate

		// send messages
		time.Sleep(time.Second * 5)
	}()
}

func main() {
	// create n clients
	clients := make([]Client, 20)

	for i := range clients {
		pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			log.Fatal(err)
		}

		pkh := sha256.New()
		pkh.Write([]byte{0x04}) // key is in uncompressed form (for now)
		pkh.Write(pk.X.Bytes())
		pkh.Write(pk.Y.Bytes())
		hash := pkh.Sum(nil)
		clients[i] = Client{
			PrivateKey:    pk,
			PublicKeyHash: base64.RawURLEncoding.EncodeToString(hash),
		}
	}

	// pair clients up
	for i := range clients {
		var other Client
		if i%2 == 0 {
			other = clients[i+1]
		} else {
			other = clients[i-1]
		}
		clients[i].Contact = Contact{
			PublicKey:     &other.PrivateKey.PublicKey,
			PublicKeyHash: other.PublicKeyHash,
		}
	}

	for i, c := range clients {
		// print out client pairs
		fmt.Println(i, c.PublicKeyHash, "paired with", c.Contact.PublicKeyHash)
		if i%2 != 0 {
			fmt.Println("--------------------------")
		}
		// start each client
		c.Start()
	}
}
