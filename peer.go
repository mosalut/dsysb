package main

import (
	"sync"
	"encoding/json"
	"net"
	"net/http"
	"time"
	"fmt"

	"github.com/mosalut/q2p"
)

const (
	p2p_post_index_event = iota
	p2p_transport_sendrawtransaction_event = 1
)

var peer *q2p.Peer_T
var receivedTransportIds = make(map[string]string)
var receivedTransportIdsMutex = &sync.RWMutex{}

func lifeCycle(peer *q2p.Peer_T, rAddr *net.UDPAddr, cycle int) {
	switch cycle {
	case q2p.JOIN:
		print(log_info, "life cycle JOIN")
		postIndex()
	case q2p.CONNECT:
		print(log_info, "life cycle CONNECT")
		postIndex()
	case q2p.CONNECTED:
		print(log_info, "life cycle CONNECTED")
	case q2p.TRANSPORT_FAILED:
		print(log_info, "life cycle TRANSPORT_FAILED")
	}
}

func addReceivedTransportId(transportId, rAddr string) {
	receivedTransportIdsMutex.Lock()
	receivedTransportIds[transportId] = rAddr
	receivedTransportIdsMutex.Unlock()
	go func(transportId string) {
		time.Sleep(30 * time.Second)
		deleteReceivedTransportId(transportId)
	} (transportId)
}

func deleteReceivedTransportId(transportId string) {
	receivedTransportIdsMutex.Lock()
	delete(receivedTransportIds, transportId)
	receivedTransportIdsMutex.Unlock()
}

func postIndex() {
	/*
	state := getState()
	index := binary.LittleEndian.Uint64(state.prevHash)
	*/
}

func transportSuccessed(peer *q2p.Peer_T, rAddr *net.UDPAddr, key string, body []byte) {
//	fmt.Println("hash key:", key)

	if len(body) < 1 {
		return
	}

	transportId := fmt.Sprintf("%056x", body[:28])
	event := uint8(body[28])
	receivedTransportIdsMutex.Lock()
	_, ok := receivedTransportIds[transportId]
	if ok {
		return
	}
	receivedTransportIdsMutex.Unlock()

	addReceivedTransportId(transportId, rAddr.String())

	switch event {
	case p2p_post_index_event:
	case p2p_transport_sendrawtransaction_event:
		tx := decodeRawTransaction(body[29:])

		err := tx.validate()
		if err != nil {
			print(log_error, err)
			return
		}
		poolMutex.Lock()
		transactionPool = append(transactionPool, tx)
		poolMutex.Unlock()
	}
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

	jsonData, err := json.Marshal(peer)
	if err != nil {
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	writeResult(w, responseResult_T{true, "ok", jsonData})
}
