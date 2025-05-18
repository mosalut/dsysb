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

func (coinbase *coinbase_T) encode() []byte {
	bs := make([]byte, coinbase_length, coinbase_length)
	copy(bs[:coinbase_amount_position], []byte(coinbase.to))
	binary.LittleEndian.PutUint64(bs[coinbase_amount_position:coinbase_nonce_position], coinbase.amount)
	binary.LittleEndian.PutUint32(bs[coinbase_nonce_position:], coinbase.nonce)

	return bs
}

func (tx *coinbase_T) encodeForPool() []byte {
	length := coinbase_length + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], coinbase_length)
	copy(bs[2:], tx.encode())

	return bs
}

func decodeCoinbase(bs []byte) *coinbase_T {
	coinbase := &coinbase_T{}
	coinbase.to = string(bs[:coinbase_amount_position])
	coinbase.amount = binary.LittleEndian.Uint64(bs[coinbase_amount_position:coinbase_nonce_position])
	coinbase.nonce = binary.LittleEndian.Uint32(bs[coinbase_nonce_position:])

	return coinbase
}

func (coinbase *coinbase_T) validate(fromP2p bool) error {
	if !fromP2p {
		return errors.New("illage type")
	}

	return nil
}

func (coinbase *coinbase_T) verifySign() bool {
	return true
}

func (coinbase *coinbase_T) count(state *state_T, c *coinbase_T, index int) error {
	_, ok := state.accounts[coinbase.to]
	if !ok {
		state.accounts[coinbase.to] = &account_T {
			coinbase.amount,
			make(map[string]uint64),
			0,
		}
	} else {
		state.accounts[coinbase.to].balance += coinbase.amount
	}

	return nil
}

func (tx *coinbase_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
		"\ttype:\tcoinbase\n" +
		"\tto: %s\n" +
		"\tamount: %d\n" +
		"\tnonce: %d", tx.hash(), tx.to, tx.nonce, tx.amount)
}
