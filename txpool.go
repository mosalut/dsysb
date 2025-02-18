package main

import (
	"sync"
	"encoding/json"
)

var transactionPool = make([]*transaction_T, 0, 511)
var poolMutex = &sync.Mutex{}

type poolCache_T struct {
	State *state_T `json:"state"`
	Transactions []*transaction_T `json:"transactions"`
}

func (cache *poolCache_T) encode() ([]byte, error) {
	bs, err := json.Marshal(cache)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func decodePoolCache(bs []byte) (*poolCache_T, error) {
	cache := &poolCache_T{}
	err := json.Unmarshal(bs, cache)
	if err != nil {
		return nil, err
	}

	return cache, nil
}
