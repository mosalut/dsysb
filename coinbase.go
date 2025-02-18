// dsysb

package main

import (
	"encoding/json"
	"fmt"
)

type coinbase_T struct {
	To string `json:"to"`
	Amount uint64 `json:"amount"`
}

func (coinbase *coinbase_T) encode() ([]byte, error) {
	bs, err := json.Marshal(coinbase)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func decodeCoinbase(bs []byte) (*coinbase_T, error) {
	coinbase := &coinbase_T{}
	err := json.Unmarshal(bs, coinbase)
	if err != nil {
		return nil, err
	}

	return coinbase, nil
}

func (coinbase *coinbase_T) String() string {
	return fmt.Sprintf(
		"\tto: %s\n" +
		"\tamount: %d", coinbase.To, coinbase.Amount)
}
