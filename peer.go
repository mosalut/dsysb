package main

import (
	"encoding/json"
	"net"
	"net/http"
	"fmt"
	"log"

	"github.com/mosalut/q2p"
)

type p2pParams_T struct {
	Key int `json:"key"`
	Data []byte `json:"data"`
}

const (
//	p2p_transport_sendrawtransaction_event = iota
	p2p_transport_sendrawtransaction_event = 1
)

var peer *q2p.Peer_T

func transportSuccessed(peer *q2p.Peer_T, rAddr *net.UDPAddr, key string, body []byte) {
//	fmt.Println("hash key:", key)

	params := &p2pParams_T{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		log.Println(err)
		return
	}

	switch params.Key {
	case p2p_transport_sendrawtransaction_event:
	//	rawtransaction := hex.EncodeToString(params.Data)

		tx := decodeRawTransaction(params.Data)

		err = tx.validate()
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
