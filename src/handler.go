package main

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"

	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/docker/distribution/notifications"
)

func statsHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	for _, tor := range torrentClient.Torrents() {
		fmt.Fprintf(w, "Torrent: %v\n", tor.Name())
		fmt.Fprintf(w, "br %v bw %v\n", tor.Stats().BytesRead, tor.Stats().BytesWritten)
	}
}

func torrentHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("unexpected request method: %v", r.Method)
		return
	}

	// Extract the content type and make sure it matches
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error parsing media type: %v, contenttype=%q", err, contentType)
		return
	}

	if mediaType != "application/octet-stream" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		log.Printf("incorrect media type: %q != %q", mediaType, "application/octet-stream")
		return
	}

	var mi metainfo.MetaInfo
	dec := bencode.NewDecoder(r.Body)
	if err := dec.Decode(&mi); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request body: %v", err)
		return
	}
	log.Printf("Got torrent, retrieving")
	seedTorrent(&mi, loadImageFromTorrent)

}

func regHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("unexpected request method: %v", r.Method)
		return
	}

	// Extract the content type and make sure it matches
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error parsing media type: %v, contenttype=%q", err, contentType)
		return
	}

	if mediaType != notifications.EventsMediaType {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		log.Printf("incorrect media type: %q != %q", mediaType, notifications.EventsMediaType)
		return
	}

	var envelope notifications.Envelope
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&envelope); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request body: %v", err)
		return
	}

	evs := envelope.Events
	for _, e := range evs {

		mu.Lock()
		if !seen[e.ID] && e.Action == "push" &&
			e.Target.MediaType == "application/vnd.docker.distribution.manifest.v2+json" {
			seen[e.ID] = true
			logEvent(e)
			//Probably need to put this in a go func, but want to make sure not rc
			downloadAndSeedImage("localhost:5000", e.Target.Repository, e.Target.Tag)
		}
		mu.Unlock()
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
	//filter on action push and mediatype manifest?
}

func hubHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("unexpected request method: %v", r.Method)
		return
	}

	// Extract the content type and make sure it matches
	contentType := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error parsing media type: %v, contenttype=%q", err, contentType)
		return
	}

	if mediaType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		log.Printf("incorrect media type: %q != %q", mediaType, "application/json")
		return
	}

	var hubEvent map[string]interface{}

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&hubEvent); err != nil {
		log.Printf("Failed to read JSON from Hub")
		panic(err)
	}
	pd := hubEvent["push_data"].(map[string]interface{})
	tag := pd["tag"].(string)
	rep := hubEvent["repository"].(map[string]interface{})
	repoName := rep["repo_name"].(string)

	if tag != "" && repoName != "" {
		go func() {
			downloadAndSeedImage("docker.io", repoName, tag)
		}()
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Recieved notification of update to %s:%s\n", repoName, tag)

}

/*
 * EOF!
 */
