// dsysb

package main

import (
	"sync"
	"encoding/binary"
)

var lenTypeM = map[int]uint8{
	coinbase_length:type_coinbase,
	create_asset_length:type_create,
	transfer_length:type_transfer,
	exchange_length:type_exchange}

var typeLenM = map[uint8]int{
	type_coinbase:coinbase_length,
	type_create:create_asset_length,
	type_transfer:transfer_length,
	type_exchange:exchange_length}

var poolMutex = &sync.Mutex{}

var signatures = make([]string, 0, 511)

type txPool_T []transaction_I

func (pool txPool_T) encode() []byte {
	bs := []byte{}

	for _, transaction := range pool {
		rawTransaction := transaction.encode()
		bs = append(bs, byte(lenTypeM[len(rawTransaction)]))
		bs = append(bs, rawTransaction...)
	}

	return bs
}

func decodeTxPool(bs []byte) txPool_T {
	bsLen := len(bs)
	var start int
	var end int
	var length int
	pool := make(txPool_T, 0, 512)
	end = len(bs)

	var typ uint8
	for start < bsLen {
		typ = bs[start]
		start++

		length = typeLenM[typ]
		end = start + length
		transaction := decodeRawTransaction(bs[start:end])
		pool = append(pool, transaction)

		start = end
	}

	return pool
}

var transactionPool = make(txPool_T, 0, 511)

type poolCache_T struct {
	state *state_T
	transactions txPool_T
}

func (cache *poolCache_T) encode() []byte {
	bs := cache.state.encode()
	stateLength := len(bs)

	rawTransactions := cache.transactions.encode()
	bs = append(bs, rawTransactions...)
	bs = append(bs, []byte{0, 0, 0, 0}...)
	binary.LittleEndian.PutUint32(bs[len(bs) - 4:], uint32(stateLength))

	return bs
}

func decodePoolCache(bs []byte) *poolCache_T {
	cache := &poolCache_T{}

	length := len(bs)
	start := length - 4
	stateLength := int(binary.LittleEndian.Uint32(bs[start:]))

	cache.state = decodeState(bs[:stateLength])

	start = stateLength
	cache.transactions = decodeTxPool(bs[start:length - 4])

	return cache
}

/*
func (cache *poolCache_T) count() {
	for k, transaction := range cache.transactions {
		transaction.count(cache, k)
	}
}

func deleteFromCacheTransactions(cache *poolCache_T, k int) {
	if len(cache.transactions) - 1 == k {
		cache.transactions = cache.transactions[:k]
	} else {
		cache.transactions = append(cache.transactions[:k], cache.transactions[k + 1:]...)
	}
}
*/
