package main

import (
	"net/http"
	"fmt"
	"log"

	"github.com/mosalut/q2p"
)

var peer *q2p.Peer_T

func callback(data []byte) {
	fmt.Println(data)
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
