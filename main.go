package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"
)

type Contact struct {
	PublicKey     *ecdsa.PublicKey
	PublicKeyHash string
}

func InitClients(clients []Client) []Client {
	// init keys
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

	// pair them up
	for i := range clients {
		var other Client
		if i%2 == 0 {
			other = clients[i+1]
		} else {
			other = clients[i-1]
			fmt.Println(clients[i].PublicKeyHash, "paired with", other.PublicKeyHash)
		}
		clients[i].Contact = Contact{
			PublicKey:     &other.PrivateKey.PublicKey,
			PublicKeyHash: other.PublicKeyHash,
		}
	}
	return clients
}

func main() {
	opt := GetOptions()
	clients := InitClients(make([]Client, opt.ClientCount))
	wg := sync.WaitGroup{}
	tStart := time.Now()
	for i := range clients {
		// start each client
		wg.Add(1)
		go clients[i].Start(&wg, opt)
		// sleep to prevent all clients connecting at once
		// this is due to an auth issue with go-npchat, which must be solved
		time.Sleep(time.Millisecond * 500)
	}
	wg.Wait()
	tEnd := time.Now()
	duration := tEnd.Sub(tStart)
	log.Println("ran in", duration.Seconds(), "seconds")
}
