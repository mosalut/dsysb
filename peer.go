package main

import (
	"math/big"
//	"crypto/elliptic"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"net"
	"net/http"
	"fmt"
	"log"
	"reflect"

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
	fmt.Println("hash key:", key)

	params := &p2pParams_T{}
	err := json.Unmarshal(body, &params)
	if err != nil {
		log.Println(err)
		return
	}

	switch params.Key {
	case p2p_transport_sendrawtransaction_event:
		rawtransaction := hex.EncodeToString(params.Data)

		transaction := transaction_T{}
		err := json.Unmarshal(params.Data, &transaction)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Println(rawtransaction)
		fmt.Println(reflect.TypeOf(transaction.PublicKey))

		publicKey := ecdsa.PublicKey{transaction.PublicKey.Curve, transaction.PublicKey.X, transaction.PublicKey.Y}
		fmt.Println(publicKey)

		ok := ecdsa.Verify(&publicKey, transaction.Txid, big.NewInt(0).SetBytes(transaction.Signature[:32]), big.NewInt(0).SetBytes(transaction.Signature[32:]))
		fmt.Println(ok)

		if ok {
			transactionPool = append(transactionPool, &transaction)
			fmt.Println(transaction)
		}
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
