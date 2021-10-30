package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const MESSAGE = `{"t":1635632889349,"iv":"PCRrbXH4rPTnO5zrTObS-0c2O-ug9WfrPdeTlCFVF2s","m":"jWMzGruE-M_oz1FrMUrkiLPTztq3UbN3xzt1rWCmOSR1OBhxx35U9yMkUNzwvG4QUSpl4tC2HR405VOg7o3LphSUC1cHAWjTU9EqRXzn0mJAQwLW0mJF80HDQi2Kw9hAuMIaFJ-L9YoKSSpdnv0s8_XFOiBJ1ARoriqzjHMBY42wGBYl1V3943N-L9AJrSU","f":"WoQSzF_Hp_sC_nv03Wv6HofbLM1iY6taN9lZ2NZyrZg","h":"agWHT5EX2-z8oSBJhkd0LFvPQru5CRKTXyPogeHffdI","p":"a-PEeDIjSAWdgESlKKo03pErIVzR9187pha_kdNTRNw","s":"Ii8BSjLXt9uLOv0jeifxU-sBn3LkIRSbyUBYu05Hz9QhHfUOHd6Ruuhf4Ke1FnUl-WxqPUXbUpwdFkT72kjzkQ"}`

type Client struct {
	PrivateKey    *ecdsa.PrivateKey
	PublicKeyHash string
	Contact       Contact
}

func (c *Client) Start(wg *sync.WaitGroup, opt Options) {
	defer wg.Done()
	// get websocket
	url := fmt.Sprintf("ws://%s/%s", opt.Hostname, c.PublicKeyHash)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to Websocket Server:", err)
	}
	defer conn.Close()

	gc := GetMessage{Get: "challenge"}
	gcJson, _ := json.Marshal(gc)
	conn.WriteMessage(websocket.TextMessage, gcJson)

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error in receive", err)
		return
	}
	var cm ChallengeMessage
	err = json.Unmarshal(msg, &cm)
	if err == nil && cm.Challenge.Txt != "" {
		// decode challenge
		txt, err := base64.RawURLEncoding.DecodeString(cm.Challenge.Txt)
		if err != nil {
			log.Fatal(err)
		}

		// hash because we're using ES256
		ch := sha256.New()
		ch.Write(txt)
		cHash := ch.Sum(nil)

		prng := rand.Reader
		r, s, err := ecdsa.Sign(prng, c.PrivateKey, cHash)
		if err != nil {
			log.Fatal(err)
		}
		sigBytes := []byte{}
		sigBytes = append(sigBytes, r.Bytes()...)
		sigBytes = append(sigBytes, s.Bytes()...)
		sigStr := base64.RawURLEncoding.EncodeToString(sigBytes)
		pubKeyBytes := []byte{}
		pubKeyBytes = append(pubKeyBytes, 0x04) // flag that key is uncompressed
		pubKeyBytes = append(pubKeyBytes, c.PrivateKey.X.Bytes()...)
		pubKeyBytes = append(pubKeyBytes, c.PrivateKey.Y.Bytes()...)
		pubKeyStr := base64.RawURLEncoding.EncodeToString(pubKeyBytes)
		cr, _ := json.Marshal(ChallengeResponse{
			Solution:  sigStr,
			PublicKey: pubKeyStr,
			Challenge: Challenge{
				Txt: cm.Challenge.Txt,
				Sig: cm.Challenge.Sig,
			},
		})
		conn.WriteMessage(websocket.TextMessage, cr)

	} else {
		log.Println(c.PublicKeyHash, string(msg))
	}

	// send messages
	for i := 0; i < opt.MessageCount; i++ {
		url := fmt.Sprintf("http://%s/%s", opt.Hostname, c.Contact.PublicKeyHash)
		r := strings.NewReader(MESSAGE)
		resp, err := http.Post(url, "text/plain", r)
		resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%s %v\n", c.PublicKeyHash, resp.StatusCode)
		_, _, err = conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		if opt.MessageWait > 0 {
			time.Sleep(time.Duration(opt.MessageWait) * time.Millisecond)
		}
	}
	err = conn.Close()
	if err != nil {
		log.Println(err)
	}
}
