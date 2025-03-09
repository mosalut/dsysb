// dsysb

package main

import (
	"crypto/sha256"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"errors"
	"fmt"
)

const (
	exchange_length = 480
)

type exchange_T [2]*transfer_T

func (ex *exchange_T) hash() [32]byte {
	length := transfer_signer_position * 2
	bs := make([]byte, length, length)
	copy(bs[:transfer_signer_position], ex[0].encodeWithoutSigner())
	copy(bs[transfer_signer_position:], ex[0].encodeWithoutSigner())
	return sha256.Sum256(ex.encode())
}

func (ex *exchange_T) getType() uint8 {
	return type_exchange
}

func (ex *exchange_T)encode() []byte {
	bs := make([]byte, exchange_length, exchange_length)
	copy(bs[:transfer_length], ex[0].encode())
	copy(bs[transfer_length:], ex[1].encode())

	return bs
}

func decodeExchange(bs []byte) *exchange_T {
	ex := &exchange_T{}
	ex[0] = decodeTransfer(bs[:transfer_length])
	ex[1] = decodeTransfer(bs[transfer_length:])

	return ex
}

func (ex *exchange_T) validate() error {
	if ex[0].from != ex[1].to || ex[0].to != ex[1].from {
		return errors.New("Exchange address not match")
	}

	state := getState()

	poolMutex.Lock()
	defer poolMutex.Unlock()

	for _, transfer := range ex {
		if transfer.from == transfer.to {
			return errors.New("Exchange to self is not allowed")
		}

		assetId := fmt.Sprintf("%064x", transfer.assetId)

		if assetId != dsysbId {
			_, ok := state.assets[assetId]
			if !ok {
				return errors.New("There's not the asset id: " + assetId)
			}
		}

		// proccess replay attack
		for _, signature := range signatures {
			s := fmt.Sprintf("%0128x", transfer.signer.signature)
			if s == signature {
				return errors.New(fmt.Sprintf("%064x", ex.hash()) + " replay: " + s)
			}
			signatures = append(signatures, s)
		}

		var nonce uint32
		account, ok := state.accounts[transfer.from]
		if !ok {
			return errors.New("There's not the account id")
		}

		nonce = account.nonce
	//	fmt.Println("nonces:", transfer.Nonce, nonce)
		if transfer.nonce - nonce != 1 {
			return errors.New("The nonces are not match")
		}

		ok = transfer.verifySign()
		if !ok {
			return errors.New("Invalid signature")
		}
	}

	return nil
}

func (ex *exchange_T) verifySign() bool {
	for _, transfer := range ex {
		publicKey := ecdsa.PublicKey{elliptic.P256(), transfer.signer.x, transfer.signer.y}
		fmt.Println(publicKey)
		txid := ex.hash()
		ok := ecdsa.Verify(&publicKey, txid[:], big.NewInt(0).SetBytes(transfer.signer.signature[:32]), big.NewInt(0).SetBytes(transfer.signer.signature[32:]))
		if !ok {
			return false
		}
	}

	return true
}

func (ex *exchange_T) count(cache *poolCache_T, index int) {
	for _, transfer := range ex {
		accountFrom, ok := cache.state.accounts[transfer.from]
		if !ok {
			print(log_warning, "Transfer from is empty address")
			deleteFromCacheTransactions(cache, index)
			return
		}

		accountTo, ok := cache.state.accounts[transfer.to]
		if !ok {
			cache.state.accounts[transfer.to] = &account_T{}
			accountTo = cache.state.accounts[transfer.to]
			accountTo.assets = make(map[string]uint64)
		}

		id := fmt.Sprintf("%064x", transfer.assetId)

		if id == dsysbId {
			if accountFrom.balance < transfer.amount {
				print(log_warning, "not enough minerals")
				deleteFromCacheTransactions(cache, index)
				return
			}

			accountFrom.balance, accountTo.balance = accountFrom.balance - transfer.amount, accountTo.balance + transfer.amount
		} else {
			balance, ok := accountFrom.assets[id]
			if !ok {
				print(log_warning, "There is not this asset")
				deleteFromCacheTransactions(cache, index)
				return
			}

			if balance < transfer.amount {
				print(log_warning, "not enough minerals")
				deleteFromCacheTransactions(cache, index)
				return
			}

			_, ok = accountTo.assets[id]
			if !ok {
				accountTo.assets[id] = 0
			}
			accountFrom.assets[id], accountTo.assets[id] = accountFrom.assets[id] - transfer.amount, accountTo.assets[id] + transfer.amount
		}
	}
}

func (ex *exchange_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
		"\tfrom: %s\n" +
		"\tto: %s\n" +
		"\tamount: %d\n" +
		"\tasset id: %064x\n" +
		"\tnonce: %d\n" +
		"%s",
		"\tfrom: %s\n" +
		"\tto: %s\n" +
		"\tamount: %d\n" +
		"\tasset id: %064x\n" +
		"\tnonce: %d\n" +
		"%s\n\n", ex.hash(), ex[0].from, ex[0].to, ex[0].amount, ex[0].assetId, ex[0].nonce, ex[0].signer, ex[1].from, ex[1].to, ex[1].amount, ex[1].assetId, ex[1].nonce, ex[1].signer)
}
