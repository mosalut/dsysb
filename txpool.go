// dsysb

package main

import (
	"sync"
	"encoding/binary"
	"encoding/hex"
	"net/http"
)

type txPool_T []transaction_I

func (pool txPool_T) encode() []byte {
	bs := []byte{}

	for _, transaction := range pool {
		rawTransactionForPool := transaction.encodeForPool()
		bs = append(bs, rawTransactionForPool...)
	}

	return bs
}

func decodeTxPool(bs []byte) txPool_T {
	bsLen := len(bs)
	var start int
	var end int
	var length int
	pool := make(txPool_T, 0, 512)

	for start < bsLen {
		end = start + 2
		length = int(binary.LittleEndian.Uint16(bs[start:end]))
		start = end
		end = start + length

		transaction, _ := decodeRawTransaction(bs[start:end])
		pool = append(pool, transaction)

		start = end
	}

	return pool
}

func (pool *txPool_T) order(transaction transaction_I) {
	for k, tx := range *pool {
		if transaction.getBytePrice() > tx.getBytePrice() {
			*pool = append(*pool, nil)
			copy((*pool)[k + 1:], (*pool)[k:])
			(*pool)[k] = transaction
			return
		}
	}

	*pool = append(*pool, transaction)
}

var transactionPool = make(txPool_T, 0, 511)
var poolMutex = &sync.RWMutex{}

func txPoolHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

//	bs := transactionPool.encode()

	bsLength := len(transactionPool) * 32
	bs := make([]byte, bsLength, bsLength)
	var start int
	end := start + 32
	for _, tx := range transactionPool {
		h := tx.hash()
		copy(bs[start:end], h[:])

		start = end
		end = start + 32
	}
	writeResult(w, responseResult_T{true, "ok", bs})
}

func txInPoolHandler(w http.ResponseWriter, req *http.Request) {
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

	for _, tx := range transactionPool {
		h := tx.hash()
		if txid == hex.EncodeToString(h[:]) {
			writeResult(w, responseResult_T{true, "ok", tx.encode()})
			return
		}

	}

	writeResult(w, responseResult_T{false, "Not found the txid in transaction pool", nil})
}

func deleteFromTransactionPool(hashStr string) {
	poolMutex.Lock()
	defer poolMutex.Unlock()

	for k, tx := range transactionPool {
		h := tx.hash()
		if hex.EncodeToString(h[:]) != hashStr {
			continue
		}

	//	fmt.Println("delete from transaction pool:", hashStr)
		if len(transactionPool) - 1 == k {
			transactionPool = transactionPool[:k]
		} else {
			transactionPool = append(transactionPool[:k], transactionPool[k + 1:]...)
		}

		break
	}
}
