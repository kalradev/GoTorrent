package torrent

import (
	"bytes"
	"crypto/sha1"
	"log"
	"net/url"
	"strconv"

	"github.com/jackpal/bencode-go"
)

type file struct {
	Length int    `bencode:"length"`
	Path   string `bencode:"path"`
}

type info struct {
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Length      int    `bencode:"length,omitempty"`
	Files       []file `bencode:"files,omitempty"`
}

type Torrent struct {
	Announce string `bencode:"announce"`
	Info     info   `bencode:"info"`
}

func (t *Torrent) BuildURL(peerID [20]byte, port uint16) string {

	log.Println("Building url")
	u, err := url.Parse(t.Announce)
	if err != nil {
		log.Fatalf("There is some problem with URL => %s", err)
	}

	// setting query parameters
	q := url.Values{
		"info_hash":  []string{t.Infohash()},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Info.Length)},
	}

	u.RawQuery = q.Encode()

	log.Println("URL built")
	return u.String()
}

func (t *Torrent) Infohash() string {
	log.Println("Computing infohash")
	return t.Info.hash()
}

func (i *info) hash() string {
	var buf bytes.Buffer
	err := bencode.Marshal(&buf, *i)
	if err != nil {
		log.Fatalln("Not able to compute infohash")
	}
	h := sha1.Sum(buf.Bytes())
	return string(h[:])
}
