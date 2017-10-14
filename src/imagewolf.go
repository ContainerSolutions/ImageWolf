package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	// Bp Declare that the program has started
	log.Println("Starting up")
	// lookupHost get LOOKUP_HOST env
	lookupHost = os.Getenv("LOOKUP_HOST")

	if lookupHost == "" {
		lookupHost = "tasks.imagewolf"
	}

	log.Printf("LOOKUP_HOST set to %s", lookupHost)

	torrentClient, _ = torrentInit()

	log.Println("Starting up")
	getMyIps()
	getPeers()

	// Start the router
	router := newRouter()
	//Registry expects to find us on port 8000
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", apiPort), router))
}

/*
 * EOF!
 */
