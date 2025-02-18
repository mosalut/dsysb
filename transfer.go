// dsysb

package main

import (
	"encoding/json"
	"fmt"
)

type transfer_T struct {
	From string `json:"from"`
	To string `json:"to"`
	Amount uint64 `json:"amount"`
	AssetId [32]byte `json:"assetId"`
	Nonce uint32 `json:"nonce"`
	Signer *signer_T `json:"signer"`
}

func (transfer *transfer_T) encode() ([]byte, error) {
	bs, err := json.Marshal(transfer)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func decodeTransfer(bs []byte) (*transfer_T, error) {
	transfer := &transfer_T{}
	err := json.Unmarshal(bs, transfer)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (transfer *transfer_T) String() string {
	return fmt.Sprintf(
		"\tfrom: %s\n" +
		"\tto: %s\n" +
		"\tamount: %d\n" +
		"\tasset id: %064x\n" +
		"\tnonce: %d\n" +
		"\t%s",
		transfer.From, transfer.To, transfer.Amount, transfer.AssetId, transfer.Nonce, transfer.Signer)
}
