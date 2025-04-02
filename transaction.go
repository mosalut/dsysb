// dsysb

package main

import (
	"strconv"
	"encoding/hex"
	"io"
	"net/http"
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
	validate(bool) error
	verifySign() bool
	countOnNewBlock(*state_T) error
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
	txid := transaction.hash()

	err := transaction.validate(false)
	if err != nil {
		return err
	}

	poolMutex.Lock()
	transactionPool = append(transactionPool, transaction)
	poolMutex.Unlock()

	broadcast(p2p_transport_sendrawtransaction_event, txid[:])

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
		if txid == fmt.Sprintf("%064x", tx.hash()) {
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
			if txid == fmt.Sprintf("%064x", tx.hash()) {
				writeResult(w, responseResult_T{true, "ok", tx.encode()})
				return
			}
		}
	}

	writeResult(w, responseResult_T{false, "Not found the txid in the last " + n + " blocks", nil})
}

func poolToCache() (*poolCache_T, error) {
	state, err := getState()
	if err != nil {
		return nil, err
	}

	var prevHash [36]byte
	/* keepit */
	block, err := getHashBlock()
	if err == errZeroBlock {
		print(log_warning, err)
	} else if err != nil {
		return nil, err
	} else {
		prevHash = [36]byte(block.head.hash)
	}

	if len(transactionPool) <= 511 {
		return &poolCache_T {
			prevHash,
			state,
			transactionPool,
		}, nil
	}

	return &poolCache_T {
		prevHash,
		state,
		transactionPool[:511],
	}, nil
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

	bs := transactionPool.encode()

	writeResult(w, responseResult_T{true, "ok", bs})
}
