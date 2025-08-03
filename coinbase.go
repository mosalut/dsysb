// dsysb

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	coinbase_length = 47
	coinbase_to_position = 1
	coinbase_amount_position = 35
	coinbase_nonce_position = 43

	three_year_blocks = 157680
)

type coinbase_T struct {
	to string
	amount uint64
	nonce uint32
}

func (coinbase *coinbase_T) hash() [32]byte {
	return sha256.Sum256(coinbase.encode())
}

func (coinbase *coinbase_T) encode() []byte {
	bs := make([]byte, coinbase_length, coinbase_length)
	bs[0] = type_coinbase
	copy(bs[coinbase_to_position:coinbase_amount_position], []byte(coinbase.to))
	binary.LittleEndian.PutUint64(bs[coinbase_amount_position:coinbase_nonce_position], coinbase.amount)
	binary.LittleEndian.PutUint32(bs[coinbase_nonce_position:], coinbase.nonce)

	return bs
}

func (coinbase *coinbase_T) encodeForPool() []byte {
	length := coinbase_length + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], coinbase_length)
	copy(bs[2:], coinbase.encode())

	return bs
}

func decodeCoinbase(bs []byte) *coinbase_T {
	coinbase := &coinbase_T{}
	coinbase.to = string(bs[coinbase_to_position:coinbase_amount_position])
	coinbase.amount = binary.LittleEndian.Uint64(bs[coinbase_amount_position:coinbase_nonce_position])
	coinbase.nonce = binary.LittleEndian.Uint32(bs[coinbase_nonce_position:])

	return coinbase
}

func (coinbase *coinbase_T) rewards(index uint32) {
	n := index / three_year_blocks // The blocks in three years.
	switch n {
	case 0:
		coinbase.amount = 5e10
	case 1:
		coinbase.amount = 25e9
	case 2:
		coinbase.amount = 125e8
	default:
		coinbase.amount = 0
	}
}

func (coinbase *coinbase_T) validate(head *blockHead_T, fromP2p bool) error {
	if !fromP2p {
		return errors.New("illage type")
	}

	index := binary.LittleEndian.Uint32(head.prevHash[32:])

	var amount uint64
	n := index / three_year_blocks // The blocks in three years.

	switch n {
	case 0:
		amount = 5e10
	case 1:
		amount = 25e9
	case 2:
		amount = 125e8
	default:
		amount = 0
	}

	fmt.Println("coinbase amount:", amount)
	if amount != coinbase.amount {
		return errors.New("The rewards and block height are not match")
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

func (coinbase *coinbase_T) getBytePrice() uint32 {
	return 0
}

func (tx *coinbase_T) Map() map[string]interface{} {
	txM := make(map[string]interface{})
	h := tx.hash()
	txM["txid"] = hex.EncodeToString(h[:])
	txM["type"] = type_coinbase
	txM["to"] = tx.to
	txM["amount"] = tx.amount
	txM["nonce"] = tx.nonce

	return txM
}

func (tx *coinbase_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
		"\ttype:\tcoinbase\n" +
		"\tto: %s\n" +
		"\tamount: %d\n" +
		"\tnonce: %d", tx.hash(), tx.to, tx.amount, tx.nonce)
}
