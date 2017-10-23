package main

import (
	"fmt"
	"log"
	"net/http"
)

/*
 * Declare the routes and their handlers
 */
var routes = Routes{

	// ImageWolf registry Routes
	Route{"registry", "POST", "/registryNotifications", regHandler},
	// ImageWolf hub Routes
	Route{"hub", "POST", "/hubNotifications", hubHandler},
	// ImageWolf torrent Routes
	Route{"torrent", "POST", "/torrent", torrentHandler},
	// ImageWolf Stats Routes
	Route{"Stats", "GET", "/stats", statsHandler},
	// API healt Routes
	Route{"APIhealt", "GET", "/healt", APIhealt},
}

/*
 * APIhealt check if AuthServer is up
 */
func APIhealt(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "ImageWolf is UP!")
	log.Println("Healt Check from ", r.RemoteAddr)
}

/*
 * END OF FILE!
 */
