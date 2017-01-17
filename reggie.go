package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/docker/distribution/notifications"
)

var seen = make(map[string]bool)
var mu sync.Mutex

var peerSet = make(map[string]bool)
var peers []net.IP
var torrentPeers []torrent.Peer //Should be ptrs IMO, but underlying lib wants copies
var torrentClient *torrent.Client
var apiPort = 8000

func main() {

	/*
		//Superseded by DNS code
			for _, host := range os.Args[1:] {
				tcpaddr, err := net.ResolveTCPAddr("tcp", host)
				if err != nil {
					//retry as other container may not be running yet
					for i := 0; err != nil && i < 100; i++ {
						time.Sleep(10)
						tcpaddr, err = net.ResolveTCPAddr("tcp", host)
					}
				}
				if err != nil {
					log.Printf("expected IP Address, got %s %v\n", host, err)
				} else {
					log.Printf("tcp addr: %v\n", tcpaddr)
					peers = append(peers, tcpaddr)
					torrentPeers = append(torrentPeers, torrent.Peer{
						IP: tcpaddr.IP, Port: 6000})
				}
			}
	*/

	var clientConfig torrent.Config
	clientConfig.Seed = true
	clientConfig.DisableTrackers = true
	clientConfig.NoDHT = true
	clientConfig.ListenAddr = "0.0.0.0:6000"
	var err error
	torrentClient, err = torrent.NewClient(&clientConfig)
	if err != nil {
		log.Printf("error creating client: %s", err)
		return
	}

	http.HandleFunc("/registryNotifications", regHandler)
	http.HandleFunc("/torrent", torrentHandler)
	http.HandleFunc("/stats", statsHandler)

	log.Println("Starting up")
	getPeers()
	//Registry expects to find us on port 8000
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", apiPort), nil))
}

func statsHandler(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	for _, tor := range torrentClient.Torrents() {
		fmt.Fprintf(w, "Torrent: %v\n", tor.Name())
		fmt.Fprintf(w, "br %v bw %v\n", tor.Stats().BytesRead, tor.Stats().BytesWritten)
	}

}

func getPeers() {

	ips, err := net.LookupIP("tasks.reggie")

	if err != nil {
		log.Printf("Error looking up tasks")
		return
	}

	for c, ip := range ips {
		log.Printf("%v %v", c, ip)
		if !peerSet[ip.String()] {

			peerSet[ip.String()] = true
			peers = append(peers, ip)
			torrentPeers = append(torrentPeers, torrent.Peer{
				IP: ip, Port: 6000})
		}
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

	if mediaType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		log.Printf("incorrect media type: %q != %q", mediaType, "application/json")
		return
	}

	var mi metainfo.MetaInfo
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&mi); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Printf("error decoding request body: %v", err)
		return
	}

	//log.Printf("got some metadata %v\n", mi.)
	log.Printf("Got torrent, retrieving")
	seedTorrent(&mi, loadImageFromTorrent)
	//TODO: check status and add peers to torrent

}

func loadImageFromTorrent(t *torrent.Torrent) {

	//should be a single file
	log.Printf("Got image file: %s\n", t.Files()[0].DisplayPath())
	log.Printf("Got: %d bytes\n", t.BytesCompleted())
	log.Printf("Not got: %d bytes\n", t.BytesMissing())

	err := exec.Command("docker", "load", "-i", t.Files()[0].Path()).Run()
	if err != nil {
		log.Printf("Failed load %v\n", err)
		return
	}
	log.Println("Loaded")

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
			downloadAndSeedImage(e.Target.Repository, e.Target.Tag)
		}
		mu.Unlock()
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")

	//filter on action push and mediatype manifest?
	//ignore layers for moment
	//then pull given image
	//call ronnie / wellbourne?

}

func downloadAndSeedImage(repo string, tag string) {
	//Could directly use API rather than docker pull?
	imageName := fmt.Sprintf("%s/%s:%s", "localhost:5000", repo, tag)
	err := exec.Command("docker", "pull", imageName).Run()
	if err != nil {
		log.Printf("Failed pull %v\n", err)
		return

	}
	log.Println("Pulled")

	tmpfile, err := ioutil.TempFile("/", repo+tag)
	if err != nil {
		log.Print(err)
		return
	}
	tmpfile.Close()
	err = exec.Command("docker", "save", "-o", tmpfile.Name(), imageName).Run()
	if err != nil {
		log.Printf("Failed save %v\n", err)
		return
	}
	log.Println("Saved")
	mi := createTorrent(tmpfile)
	seedTorrent(&mi, notifyPeers)
	log.Println("Seeded")

	//log.Printf("torrent: %v\n", mi)

}

func seedTorrent(mi *metainfo.MetaInfo, cb func(*torrent.Torrent)) {

	//Hmm, will need to do this separately and send torrents dynamically
	//but for the moment...

	t, err := torrentClient.AddTorrent(mi)
	getPeers()

	t.AddPeers(torrentPeers)
	if err != nil {
		log.Printf("error adding torrent: %s", err)
		return
	}
	go func() {
		<-t.GotInfo()
		t.DownloadAll()
		for t.BytesMissing() > 0 {
			time.Sleep(1 * time.Second)
			fmt.Printf("Got: %d bytes missing %d\n", t.BytesCompleted(), t.BytesMissing())
		}
		cb(t)
	}()

}

func notifyPeers(t *torrent.Torrent) {

	getPeers()
	mi := t.Metainfo()
	data, err := json.Marshal(mi)
	if err != nil {
		log.Printf("Failed to create JSON %v\n", err)
	}
	r := bytes.NewReader(data)

	for _, ip := range peers {

		url := fmt.Sprintf("http://%s:%d/torrent", ip.String(), apiPort)
		fmt.Printf("Notifying: %s\n", url)

		//log.Printf("Sending torrent JSON: %s", string(data))
		http.Post(url, "application/json", r)
	}

}

func createTorrent(f *os.File) metainfo.MetaInfo {
	mi := metainfo.MetaInfo{}

	mi.SetDefaults()
	info := metainfo.Info{
		PieceLength: 256 * 1024,
	}

	err := info.BuildFromFilePath(f.Name())
	if err != nil {
		log.Fatal(err)
	}
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		log.Fatal(err)
	}

	return mi

}

func logEvent(e notifications.Event) {

	log.Println("EVENT")
	log.Printf("ACTION: %v\n", e.Action)
	log.Printf("ACTOR: %v\n", e.Actor)
	log.Printf("ID: %v\n", e.ID)
	log.Printf("REQUEST: %v\n", e.Request)
	log.Printf("SOURCE: %v\n", e.Source)
	log.Printf("TARGET DESCRIPTOR: %v\n", e.Target.Descriptor)
	log.Printf("TARGET FROMREPO: %v\n", e.Target.FromRepository)
	log.Printf("TARGET MEDIATYPE: %v\n", e.Target.MediaType)
	log.Printf("TARGET REPO: %v\n", e.Target.Repository)
	log.Printf("TARGET TAG: %v\n", e.Target.Tag)
	log.Printf("TARGET URLs: %v\n", e.Target.URLs)
	log.Printf("TARGET URL: %v\n", e.Target.URL)
	log.Printf("TIMESTAMP: %v\n", e.Timestamp)
}
