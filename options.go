package main

import (
	"flag"
	"log"
)

type Options struct {
	Hostname     string
	MessageCount int
	MessageWait  int // ms
	ClientCount  int
}

func GetOptions() Options {
	o := Options{}
	flag.StringVar(&o.Hostname, "h", "", "hostname, must be a string (with optional port)")
	flag.IntVar(&o.MessageCount, "m", 1000, "message count, must be an int")
	flag.IntVar(&o.MessageWait, "w", 5, "message wait (milliseconds), must be an int")
	flag.IntVar(&o.ClientCount, "c", 10, "client count, must be an int")
	flag.Parse()
	if o.Hostname == "" {
		log.Fatal("Please specify a hostname with -h [DOMAIN]:[PORT]")
	}
	return o
}
