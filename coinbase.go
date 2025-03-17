// dsysb

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	coinbase_length = 46
	coinbase_amount_position = 34
	coinbase_nonce_position = 42
)

type coinbase_T struct {
	to string
	amount uint64
	nonce uint32
}

func (tx *coinbase_T) hash() [32]byte {
	return sha256.Sum256(tx.encode())
}

func (coinbase *coinbase_T) getType() uint8 {
	return type_coinbase
}

func (coinbase *coinbase_T) encode() []byte {
	bs := make([]byte, coinbase_length, coinbase_length)
	copy(bs[:coinbase_amount_position], []byte(coinbase.to))
	binary.LittleEndian.PutUint64(bs[coinbase_amount_position:coinbase_nonce_position], coinbase.amount)
	binary.LittleEndian.PutUint32(bs[coinbase_nonce_position:], coinbase.nonce)

	return bs
}

func decodeCoinbase(bs []byte) *coinbase_T {
	coinbase := &coinbase_T{}
	coinbase.to = string(bs[:coinbase_amount_position])
	coinbase.amount = binary.LittleEndian.Uint64(bs[coinbase_amount_position:coinbase_nonce_position])
	coinbase.nonce = binary.LittleEndian.Uint32(bs[coinbase_nonce_position:])

	return coinbase
}

func (coinbase *coinbase_T) validate() error {
	return errors.New("illage type")
}

func (coinbase *coinbase_T) verifySign() bool {
	return true
}

func (coinbase *coinbase_T) count(cache *poolCache_T, index int) {
	_, ok := cache.state.accounts[coinbase.to]
	if !ok {
		cache.state.accounts[coinbase.to] = &account_T {
			coinbase.amount,
			make(map[string]uint64),
			0,
		}
	} else {
		cache.state.accounts[coinbase.to].balance += coinbase.amount
	}
}

func (tx *coinbase_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
		"\ttype:\tcoinbase\n" +
		"\tto: %s\n" +
		"\tamount: %d\n" +
		"\tnonce: %d\n\n", tx.hash(), tx.to, tx.nonce, tx.amount)
}
