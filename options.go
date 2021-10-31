package main

import (
	"flag"
	"log"
)

type Options struct {
	Host         string // host & port
	ClientCount  int
	MessageCount int
	MessageWait  int // ms
	ClientDelay  int // ms
	Repeat       int
}

func GetOptions() Options {
	o := Options{}
	flag.StringVar(&o.Host, "h", "", "host, must be a string (with optional port)")
	flag.IntVar(&o.MessageCount, "m", 100, "message count, must be an int")
	flag.IntVar(&o.MessageWait, "w", 10, "message wait (ms), must be an int")
	flag.IntVar(&o.ClientCount, "c", 100, "client count, must be an int")
	flag.IntVar(&o.ClientDelay, "d", 10, "client delay (ms), must be an int")
	flag.IntVar(&o.Repeat, "r", 1, "repeat, must be an int")
	flag.Parse()
	if o.Host == "" {
		log.Fatal("Please specify a host with -h [DOMAIN]:[PORT]")
	}
	return o
}
