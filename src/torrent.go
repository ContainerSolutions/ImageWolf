package main

import (
	"log"

	"github.com/anacrolix/torrent"
)

func torrentInit() (*torrent.Client, error) {
	var clientConfig torrent.Config
	clientConfig.DataDir = dataDir
	clientConfig.Seed = true
	clientConfig.DisableTrackers = true
	clientConfig.NoDHT = true
	clientConfig.ListenAddr = "0.0.0.0:6000"
	var err error
	torrentClient, err = torrent.NewClient(&clientConfig)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}
	return torrentClient, err
}

/*
 * EOF!
 */
