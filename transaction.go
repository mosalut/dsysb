// dsysb

package main

import (
	"strconv"
	"encoding/hex"
	"io"
	"net/http"
)

const (
	type_coinbase = iota
	type_create
	type_transfer
	type_exchange
	type_deploy
	type_call
	type_extension
)

type transaction_I interface {
	hash() [32]byte
	encode() []byte
	validate(*blockHead_T, bool) error
	verifySign() bool
	count(*state_T, *coinbase_T, int) error
	encodeForPool() []byte
	getBytePrice() uint32
	Map() map[string]interface{}
	String() string
}

func decodeRawTransaction(bs []byte) (transaction_I, error) {
	var tx transaction_I

	switch bs[0] {
	case type_coinbase:
		tx = decodeCoinbase(bs)
	case type_create:
		tx = decodeCreateAsset(bs)
	case type_transfer:
		tx = decodeTransfer(bs)
	case type_exchange:
		tx = decodeExchange(bs)
	case type_deploy:
		tx = decodeDeployTask(bs)
	case type_call:
		tx = decodeCallTask(bs)
	case type_extension:
		tx = decodeExtension(bs)
	default:
		return nil, errWrongType
	}

	return tx, nil
}

func sendRawTransaction(bs []byte) error {
	transaction, err := decodeRawTransaction(bs)
	if err != nil {
		return err
	}

	err = transaction.validate(nil, false)
	if err != nil {
		return err
	}

	poolMutex.Lock()
	transactionPool.order(transaction)
	poolMutex.Unlock()

	broadcast(p2p_transport_sendrawtransaction_event, bs)

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

func getTransactionHandler(w http.ResponseWriter, req *http.Request) {
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
	txid := values.Get("txid")
	n := values.Get("number")
	number, err := strconv.Atoi(n)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	block, err := getHashBlock()
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	for _, tx := range block.body.transactions {
		h := tx.hash()
		if txid == hex.EncodeToString(h[:]) {
			writeResult(w, responseResult_T{true, "ok", tx.encode()})
			return
		}
	}

	for i := 0; i < number && block != nil; i++ {
		if hex.EncodeToString(block.head.prevHash[:]) == genesisPrevHash {
			break
		}

		block, err = getBlock(block.head.prevHash[32:])
		if err != nil {
			writeResult(w, responseResult_T{false, err.Error(), nil})
			return
		}

		for _, tx := range block.body.transactions {
			h := tx.hash()
			if txid == hex.EncodeToString(h[:]) {
				writeResult(w, responseResult_T{true, "ok", tx.encode()})
				return
			}
		}
	}

	writeResult(w, responseResult_T{false, "Not found the txid in the last " + n + " blocks", nil})
}
