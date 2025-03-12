// dsysb

package main

import (
	"math/rand"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"net"
	"net/http"
	"time"
	"fmt"
)

const (
	type_coinbase = iota
	type_create
	type_transfer
	type_exchange
)

type transaction_I interface {
	hash() [32]byte
	getType() uint8
	encode() []byte
	validate() error
	verifySign() bool
	count(*poolCache_T, int)
	String() string
}

func decodeRawTransaction(bs []byte) transaction_I {
	length := len(bs)

	var tx transaction_I
	switch length {
	case coinbase_length:
		tx = decodeCoinbase(bs)
	case create_asset_length:
		tx = decodeCreateAsset(bs)
	case transfer_length:
		tx = decodeTransfer(bs)
	case exchange_length:
		tx = decodeExchange(bs)
	default:
		print(log_error, "Wrong type")
	}

	return tx
}

func sendRawTransaction(bs []byte) error {
	transaction := decodeRawTransaction(bs)

	err := transaction.validate()
	if err != nil {
		return err
	}

	poolMutex.Lock()
	transactionPool = append(transactionPool, transaction)
	poolMutex.Unlock()

	originLength := 17 + len(bs)
	origin := make([]byte, originLength, originLength)

	timestamp := time.Now().UnixNano()
	binary.LittleEndian.PutUint64(origin[:8], uint64(timestamp))

	rNonce := rand.Uint64()
	binary.LittleEndian.PutUint64(origin[8:16], rNonce)
	origin[16] = p2p_transport_sendrawtransaction_event

	copy(origin[17:], bs)

	hash := sha256.Sum224(origin)

	data := append(hash[:], byte(p2p_transport_sendrawtransaction_event))
	data = append(data, bs...)

	for k, _ := range seedAddrs {
		rAddr, err := net.ResolveUDPAddr("udp", k)
		if err != nil {
			continue
			print(log_error, err)
		}

		_, err = peer.Transport(rAddr, data)
		if err != nil {
			continue
			print(log_error, err)
		}
	}

	return nil
}

func sendRawTransactionHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	err = sendRawTransaction(body)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	writeResult(w, responseResult_T{true, "ok", nil})
}

func poolToCache() *poolCache_T {
	state := getState()

	if len(transactionPool) <= 511 {
		return &poolCache_T {
			state,
			transactionPool,
		}
	}

	return &poolCache_T {
		state,
		transactionPool[:511],
	}
}

func txPool(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	fmt.Println(transactionPool)
	bs := transactionPool.encode()

	writeResult(w, responseResult_T{true, "ok", bs})
}
