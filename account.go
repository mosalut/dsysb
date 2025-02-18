package main

import (
	"crypto/sha256"
	"encoding/json"
)

type account_T struct {
	Balance uint64 `json:"balance"`
	Asset map[string]uint64 `json:"asset"` // key is an asset id
	Nonce uint32 `json:"nonce"`
}

func (account *account_T)encode() []byte {
	bs, err := json.Marshal(account)
	if err != nil {
		print(log_error, err)
		return nil
	}

	return bs
}

func decodeAccount(bs []byte) *account_T {
	account := &account_T{}
	err := json.Unmarshal(bs, account)
	if err != nil {
		print(log_error, err)
		return nil
	}

	return account
}

func (account *account_T)hash() [32]byte {
	return sha256.Sum256(account.encode())
}

var accountPool = make([]*account_T, 0, 500)
