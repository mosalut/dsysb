// dsysb

package main

import (
	"encoding/json"
)

type exchange_T [2]*transfer_T

func (ex *exchange_T)encode() ([]byte, error) {
	bs, err := json.Marshal(ex)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func decodeExchange(bs []byte) (*exchange_T, error) {
	ex := &exchange_T{}
	err := json.Unmarshal(bs, ex)
	if err != nil {
		return nil, err
	}

	return ex, nil
}

func (ex *exchange_T) String() string {
	return ex[0].String() + ex[1].String()
}

func (ex *exchange_T) check() error {
	// TODO
	return nil
}
