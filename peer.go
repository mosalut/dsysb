package main

import (
	"math/rand"
	"sync"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"net"
	"net/http"
	"time"
	"fmt"

	"github.com/mosalut/q2p"
)

const (
	p2p_post_index_event = iota
	p2p_transport_sendrawtransaction_event
	p2p_debug
)

var peer *q2p.Peer_T
var receivedTransportIds = make(map[string]string)
var receivedTransportIdsMutex = &sync.RWMutex{}

func lifeCycle(peer *q2p.Peer_T, rAddr *net.UDPAddr, cycle int) {
	switch cycle {
	case q2p.JOIN:
		print(log_info, "life cycle JOIN")
	//	postLastHash(peer, rAddr)
		postDebug()
	case q2p.CONNECT:
		print(log_info, "life cycle CONNECT")
	//	postLastHash(peer, rAddr)
	case q2p.CONNECTED:
		print(log_info, "life cycle CONNECTED")
	case q2p.TRANSPORT_FAILED:
		print(log_info, "life cycle TRANSPORT_FAILED")
	}
}

func makePostHash(bs []byte) [28]byte {
	originLength := 17 + len(bs)
	origin := make([]byte, originLength, originLength)

	timestamp := time.Now().UnixNano()
	binary.LittleEndian.PutUint64(origin[:8], uint64(timestamp))

	rNonce := rand.Uint64()
	binary.LittleEndian.PutUint64(origin[8:16], rNonce)
	origin[16] = p2p_transport_sendrawtransaction_event

	copy(origin[17:], bs)

	return sha256.Sum224(origin)
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

func postLastHash(rAddr *net.UDPAddr) {
	state := getState()
	prevHash := state.prevHash
	print(log_debug, "prevHash:", prevHash)

	hash := makePostHash(prevHash[:])

	data := append(hash[:], byte(p2p_post_index_event))
	data = append(data, prevHash[:]...)

	_, err := peer.Transport(rAddr, data)
	if err != nil {
		print(log_error, err)
	}
}

func postDebug() {
	hi := []byte("hihihihi")
	hash := makePostHash(hi)
	data := append(hash[:], byte(p2p_debug))
	data = append(data, hi...)

	broadcast(string(hash[:]), data)
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
		print(log_debug, "prevHash:", body[29:])
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
	case p2p_debug:
		print(log_debug, "postId:", fmt.Sprintf("%056x", body[:29]))
		print(log_debug, "hi:", string(body[29:]))
		broadcast(string(body[:29]), body)
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

func broadcast(transportId string, data []byte) {
	addReceivedTransportId(transportId, peer.Conn.LocalAddr().String())
	for k, _ := range seedAddrs {
		rAddr, err := net.ResolveUDPAddr("udp", k)
		if err != nil {
			print(log_error, err)
			continue
		}

		_, err = peer.Transport(rAddr, data)
		if err != nil {
			print(log_error, err)
			continue
		}
	}
}
