package main

import (
	"net"
	"net/http"
	"fmt"
	"log"

	"github.com/mosalut/q2p"
)

var peer *q2p.Peer_T

func transportSuccessed(peer *q2p.Peer_T, rAddr *net.UDPAddr, key string, data []byte) {
	fmt.Println(key, data)
}

func transportFailed(peer *q2p.Peer_T, rAddr *net.UDPAddr, key string, syns []uint32) {
	fmt.Println(key, syns)
}

func peerHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	log.Println(peer)
	writeResult(w, responseResult_T{true, "ok", *peer})
}
