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
	p2p_add_block_event
	p2p_get_block_hashes_event
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
	//	postDebug()
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
	originLength := 8 + len(bs)
	origin := make([]byte, originLength, originLength)

	rNonce := rand.Uint64()
	binary.LittleEndian.PutUint64(origin[:8], rNonce)

	copy(origin[8:], bs)

	return sha256.Sum224(origin)
}

func addReceivedTransportId(postId, rAddr string) {
	receivedTransportIdsMutex.Lock()
	receivedTransportIds[postId] = rAddr
	receivedTransportIdsMutex.Unlock()
	go func(postId string) {
		time.Sleep(30 * time.Second)
		deleteReceivedTransportId(postId)
	} (postId)
}

func deleteReceivedTransportId(postId string) {
	receivedTransportIdsMutex.Lock()
	delete(receivedTransportIds, postId)
	receivedTransportIdsMutex.Unlock()
}

func postDebug() {
	hi := []byte("hihihihi")

	broadcast(p2p_debug, hi)
}

func transportSuccessed(peer *q2p.Peer_T, rAddr *net.UDPAddr, key string, body []byte) {
	fmt.Println("hash key:", key)

	if len(body) < 1 {
		return
	}

	postId := fmt.Sprintf("%056x", body[:28])
	fmt.Println(postId)
	fmt.Println(receivedTransportIds)
	event := uint8(body[28])
	receivedTransportIdsMutex.Lock()
	_, ok := receivedTransportIds[postId]
	if ok {
		return
	}
	receivedTransportIdsMutex.Unlock()

	addReceivedTransportId(postId, rAddr.String())

	fmt.Println("inner event:", event)
	switch event {
	case p2p_transport_sendrawtransaction_event:
		tx := decodeRawTransaction(body[29:])

		err := tx.validate(true)
		if err != nil {
			print(log_error, err)
			return
		}
		poolMutex.Lock()
		transactionPool = append(transactionPool, tx)
		poolMutex.Unlock()
	case p2p_add_block_event:
		block := decodeBlock(body[29:])

		blockHash32 := fmt.Sprintf("%064x", block.head.hashing())
		print(log_debug, "hash32:")
		print(log_debug, blockHash32)
		print(log_debug, fmt.Sprintf("%072x", block.head.hash))
		if blockHash32 != fmt.Sprintf("%064x", block.head.hash[:32]) {
			print(log_warning, "p2p_add_block_event: Block hash32 not match")
			return
		}

		blockIndex := binary.LittleEndian.Uint32(block.head.hash[32:])
		blockPrevIndex := binary.LittleEndian.Uint32(block.head.prevHash[32:])

		if blockIndex - blockPrevIndex != 1 {
			print(log_warning, "p2p_add_block_event: Block index not match")
			return
		}

		state := getState()
		statePrevHash := fmt.Sprintf("%072x", state.prevHash)
		blockPrevHash := fmt.Sprintf("%072x", block.head.prevHash)
		print(log_debug, "statePrevHash:", statePrevHash)
		print(log_debug, "blockPrevHash:", blockPrevHash)
		if statePrevHash == blockPrevHash {
			err := block.Append(state)
			if err != nil {
				print(log_error, err)
				return
			}
		}

		statePrevIndex := binary.LittleEndian.Uint32(state.prevHash[32:])
		if statePrevIndex > blockPrevIndex {
			print(log_warning, "p2p_add_block_event: Get a lower block.")
			return
		}

		print(log_debug, "Block synchronization.")
		err := transport(rAddr, p2p_get_block_hashes_event, state.prevHash[:])
		if err != nil {
			print(log_error, err)
			return
		}
	case p2p_get_block_hashes_event:
		print(log_info, "p2p_get_block_hashes_event:")
		print(log_debug, "state prev hash", fmt.Sprintf("%x\n", body[29:]))
	case p2p_debug:
		postId := fmt.Sprintf("%056x", body[:28])
		print(log_debug, "postId:", postId)
		print(log_debug, "hi:", string(body[29:]))
		broadcast(p2p_debug, body)
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

func broadcast(event uint8, data []byte) {
	hash := makePostHash(data)
	bs := append(hash[:], byte(event))
	bs = append(bs, data...)
	postId := fmt.Sprintf("%056x", hash[:])

	fmt.Println("remote seeds:", peer.RemoteSeeds)
	addReceivedTransportId(postId, peer.Conn.LocalAddr().String())
	for k, _ := range peer.RemoteSeeds {
		rAddr, err := net.ResolveUDPAddr("udp", k)
		if err != nil {
			print(log_error, err)
			continue
		}

		_, err = peer.Transport(rAddr, bs)
		if err != nil {
			print(log_error, err)
			continue
		}
	}
}

func transport(rAddr *net.UDPAddr, event uint8, data []byte) error {
	hash := makePostHash(data)
	bs := append(hash[:], byte(event))
	bs = append(bs, data...)
	postId := fmt.Sprintf("%056x", hash[:])
	addReceivedTransportId(postId, peer.Conn.LocalAddr().String())

	_, err := peer.Transport(rAddr, bs)
	if err != nil {
		return err
	}

	return nil
}
