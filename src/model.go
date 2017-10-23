package main

import (
	"context"
	"net"
	"net/http"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/docker/docker/client"
)

var seen = make(map[string]bool)
var mu sync.Mutex

var peerSet = make(map[string]bool)
var peers []net.IP
var torrentPeers []torrent.Peer //Should be ptrs IMO, but underlying lib wants copies
var torrentClient *torrent.Client
var apiPort = 8000
var dataDir = "/data"
var myIps = make(map[string]bool)
var lookupHost string

// Docker ctx
var (
	CTX context.Context
	CLI *client.Client
)

/*
 * Type Route
 */
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}
type Routes []Route

/*
 * HttpResp Struct for json responses
 */
type HttpResp struct {
	Status      int    `json:"status"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

/*
 * END OF FILE!
 */
