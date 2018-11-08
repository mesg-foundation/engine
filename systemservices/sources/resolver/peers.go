package main

import (
	"sync"
)

var (
	peers      []string
	peersMutex sync.RWMutex
)

func addPeers(ps ...string) {
	peersMutex.Lock()
	defer peersMutex.Unlock()
	peers = append(peers, ps...)
}

func getPeers() []string {
	peersMutex.RLock()
	defer peersMutex.RUnlock()
	return peers
}
