package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/docker/distribution/notifications"
)

func downloadAndSeedImage(registry, repo, tag string) {

	err := dockerPull(registry, repo, tag)
	if err != nil {
		log.Printf("Failed pull of %s:%s %v\n", repo, tag, err)
		return
	}

	tmpfile, err := ioutil.TempFile(dataDir, strings.Replace(repo+tag, "/", "", -1))
	if err != nil {
		log.Print(err)
		return
	}
	defer tmpfile.Close()

	strm, err := dockerSave(registry, repo, tag)
	if err != nil {
		log.Print(err)
		return
	}

	_, err = io.Copy(tmpfile, strm)

	if err != nil {
		log.Printf("Failed save %v\n", err)
		return
	}

	log.Println("Saved")
	mi := createTorrent(tmpfile)
	seedTorrent(&mi, notifyPeers)
	log.Println("Seeding")

}

func seedTorrent(mi *metainfo.MetaInfo, cb func(*torrent.Torrent)) {

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
			log.Printf("Got: %d bytes missing %d\n", t.BytesCompleted(), t.BytesMissing())
		}
		cb(t)
	}()

}

func notifyPeers(t *torrent.Torrent) {

	getPeers()
	mi := t.Metainfo()
	data, err := bencode.Marshal(mi)
	if err != nil {
		log.Printf("Failed to create JSON %v\n", err)
	}

	for _, ip := range peers {

		url := fmt.Sprintf("http://%s:%d/torrent", ip.String(), apiPort)
		log.Printf("Notifying: %s\n", url)

		_, err := http.Post(url, "application/octet-stream", bytes.NewReader(data))
		if err != nil {
			log.Printf("notify responded with err %v\n", err)
		}
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

func loadImageFromTorrent(t *torrent.Torrent) {

	//should be a single file
	log.Printf("Got: %d bytes\n", t.BytesCompleted())
	log.Printf("Not got: %d bytes\n", t.BytesMissing())

	err := dockerLoad(dataDir + "/" + t.Files()[0].Path())
	if err != nil {
		log.Printf("Failed load %v\n", err)
		return
	}

	log.Println("Loaded image from torrent ", t.Files()[0].Path())

}

func getPeers() {

	ips, err := net.LookupIP(lookupHost)

	if err != nil {
		log.Printf("Error looking up hosts")
		return
	}

	for c, ip := range ips {

		if !peerSet[ip.String()] && !myIps[ip.String()] {
			log.Printf("%v %v", c, ip)

			peerSet[ip.String()] = true
			peers = append(peers, ip)
			torrentPeers = append(torrentPeers, torrent.Peer{
				IP: ip, Port: 6000})
		}
	}

}

// getMyIps get the Ip of the ImageWolf host
func getMyIps() {

	ifaces, err := net.Interfaces()
	if err != nil {
		log.Printf("Failed to inspect my network interfaces %v\n", err)
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			log.Printf("Failed to inspect network addresses %v\n", err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			log.Printf("Found my IP %s\n", ip.String())
			myIps[ip.String()] = true
		}
	}
}

/*
 * END OF FILE!
 */
