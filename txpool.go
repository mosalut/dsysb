// dsysb

package main

import (
	"sync"
	"encoding/binary"
	"encoding/hex"
	"net/http"
)

var txIds = make([]string, 0, 511)
var txIdsMutex = &sync.RWMutex{}

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

var transactionPool = make(txPool_T, 0, 511)
var poolMutex = &sync.RWMutex{}

type poolCache_T struct {
	prevHash [36]byte
	bits [4]byte
	state *state_T
	transactions txPool_T
}

func (cache *poolCache_T) encode() []byte {
	bs := cache.prevHash[:]
	bs = append(bs, cache.bits[:]...)
	bs = append(bs, cache.state.encode()...)
	transactionsPosition := len(bs)

	rawTransactions := cache.transactions.encode()
	bs = append(bs, rawTransactions...)
	bs = append(bs, []byte{0, 0, 0, 0}...)
	binary.LittleEndian.PutUint32(bs[len(bs) - 4:], uint32(transactionsPosition))

	return bs
}

func decodePoolCache(bs []byte) *poolCache_T {
	cache := &poolCache_T{}

	length := len(bs)
	cache.prevHash = [36]byte(bs[:36])
	cache.bits = [4]byte(bs[36:40])

	start := length - 4
	transactionsPosition := int(binary.LittleEndian.Uint32(bs[start:]))

	cache.state = decodeState(bs[40:transactionsPosition])

	start = transactionsPosition
	cache.transactions = decodeTxPool(bs[start:length - 4])

	return cache
}

/* keepfunc */
func poolToCache() (*poolCache_T, error) {
	var prevHash [36]byte

	// keepit
	block, err := getHashBlock()
	if err == errZeroBlock {
		print(log_warning, err)
	} else if err != nil {
		return nil, err
	} else {
		prevHash = block.head.hash
	}

	err = adjustTarget(block)
	if err != nil {
		print(log_error, err)
		return nil, err
	}

	if len(transactionPool) <= 511 {
		return &poolCache_T {
			prevHash,
		//	difficult_1_target, // keepit
			block.head.bits,
		//	firstState, // keepit
			block.state,
			transactionPool,
		}, nil
	}

	// keepit
	return &poolCache_T {
		prevHash,
		block.head.bits,
		block.state,
		transactionPool[:511],
	}, nil
}

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

	bs := transactionPool.encode()

	writeResult(w, responseResult_T{true, "ok", bs})
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
