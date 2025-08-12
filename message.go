package main

import (
	"net"
	"net/http"
	"fmt"
)

func sendMessageHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	values := req.URL.Query()
	addr := values.Get("address")
	message := values.Get("message")

	sendMessage(addr, message)
}


func sendMessage(addr, message string) error {
	rAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	print(log_debug, message, rAddr)
	hash, err := peer.Transport(rAddr, []byte(message))
	if err != nil {
		return err
	}

	fmt.Println(hash, "sent")

	return nil
}

func broadcastHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	values := req.URL.Query()
	message := values.Get("message")

	broadcast(p2p_debug, []byte(message))

	writeResult(w, responseResult_T{true, "ok", nil})
}
