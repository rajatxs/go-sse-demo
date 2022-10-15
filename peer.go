package main

type Peer struct {
	Id string
	Ch chan string
}

var activePeers map[string]*Peer = map[string]*Peer{}
