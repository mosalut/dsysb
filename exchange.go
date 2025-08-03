// dsysb

package main

import (
	"crypto/sha256"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"errors"
	"fmt"
)

const (
	exchange_length = 491
)

type exchange_T [2]*transfer_T

func (ex *exchange_T) hash() [32]byte {
	length := transfer_signer_position * 2 + 1
	bs := make([]byte, length, length)
	bs[0] = type_exchange
	copy(bs[1:transfer_signer_position + 1], ex[0].encodeWithoutSigner())
	copy(bs[transfer_signer_position + 1:], ex[1].encodeWithoutSigner())
	return sha256.Sum256(bs)
}

func (ex *exchange_T)encode() []byte {
	bs := make([]byte, exchange_length, exchange_length)
	bs[0] = type_exchange
	copy(bs[1:transfer_length + 1], ex[0].encode())
	copy(bs[transfer_length + 1:], ex[1].encode())

	return bs
}

func (ex *exchange_T)encodeForPool() []byte {
	length := exchange_length + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], exchange_length)
	copy(bs[2:], ex.encode())

	return bs
}

func decodeExchange(bs []byte) *exchange_T {
	ex := &exchange_T{}
	ex[0] = decodeTransfer(bs[1:transfer_length +1])
	ex[1] = decodeTransfer(bs[transfer_length + 1:])

	return ex
}

func (ex *exchange_T) validate(head *blockHead_T, fromP2p bool) error {
	if ex[0].from != ex[1].to || ex[0].to != ex[1].from {
		return errors.New("The exchange addresses are not match")
	}

	if hex.EncodeToString(ex[0].assetId[:]) == hex.EncodeToString(ex[1].assetId[:]) {
		return errors.New("The assetIds is the same one")
	}

	state, err := getState()
	if err != nil {
		return err
	}

	// replay attack
	txIdH := ex.hash()
	txId := hex.EncodeToString(txIdH[:])
	for k, tx := range transactionPool {
		h := tx.hash()
		if txId == hex.EncodeToString(h[:]) {
			if fromP2p {
			//	deleteFromTransactionPool(txId)
				poolMutex.Lock()
				if len(transactionPool) - 1 == k {
					transactionPool = transactionPool[:k]
				} else {
					transactionPool = append(transactionPool[:k], transactionPool[k + 1:]...)
				}
				poolMutex.Unlock()
				return nil
			}

			return errors.New("Replay attack: txid: " + txId)
		}
	}

	for _, transfer := range ex {
		if transfer.from == transfer.to {
			return errors.New("Exchange to self is not allowed")
		}

		assetId := fmt.Sprintf("%064x", transfer.assetId)

		if assetId != dsysbId {
			asset, ok := state.assets[assetId]
			if !ok {
				return errors.New("There's not the asset id: " + assetId)
			}

			if transfer.bytePrice < asset.price {
				return errors.New(fmt.Sprintf("The byte price should >= asset's create price: %d", asset.price))
			}
		}

		var nonce uint32
		account, ok := state.accounts[transfer.from]
		if !ok {
			return errors.New("There's not the account id")
		}

		nonce = account.nonce
		if transfer.nonce - nonce != 1 {
			return errOutOfNonce
		}
	}

	ok := ex.verifySign()
	if !ok {
		return errors.New("Invalid signature")
	}

	return nil
}

func (ex *exchange_T) verifySign() bool {
	for _, transfer := range ex {
		publicKey := ecdsa.PublicKey{elliptic.P256(), transfer.signer.x, transfer.signer.y}
		txid := ex.hash()
		ok := ecdsa.Verify(&publicKey, txid[:], big.NewInt(0).SetBytes(transfer.signer.signature[:32]), big.NewInt(0).SetBytes(transfer.signer.signature[32:]))
		if !ok {
			print(log_info, "Invalid signature")
			return false
		}
	}

	return true
}

func (ex *exchange_T) count(state *state_T, coinbase *coinbase_T, index int) error {
	for _, transfer := range ex {
		accountFrom, ok := state.accounts[transfer.from]
		if !ok {
			return errors.New("Transfer from is empty address")
		}

		accountTo, ok := state.accounts[transfer.to]
		if !ok {
			state.accounts[transfer.to] = &account_T{}
			accountTo = state.accounts[transfer.to]
			accountTo.assets = make(map[string]uint64)
		}

		id := fmt.Sprintf("%064x", transfer.assetId)

		if id == dsysbId {
			if accountFrom.balance < transfer.amount {
				return errors.New("not enough minerals")
			}

			accountFrom.balance, accountTo.balance = accountFrom.balance - transfer.amount, accountTo.balance + transfer.amount
		} else {
			balance, ok := accountFrom.assets[id]
			if !ok {
				return errors.New("There is not this asset")
			}

			if balance < transfer.amount {
				return errors.New("not enough minerals")
			}

			_, ok = accountTo.assets[id]
			if !ok {
				accountTo.assets[id] = 0
			}
			accountFrom.assets[id], accountTo.assets[id] = accountFrom.assets[id] - transfer.amount, accountTo.assets[id] + transfer.amount
		}

		accountFrom.balance -= transfer.fee()
	//	state.accounts[*address].balance += transfer.fee()
		coinbase.amount += transfer.fee()
		accountFrom.nonce = transfer.nonce
	}

	return nil
}

func (ex *exchange_T) getBytePrice() uint32 {
	return (ex[0].bytePrice + ex[1].bytePrice) / 2
}

func (ex *exchange_T) Map() map[string]interface{} {
	txM := make(map[string]interface{})

	h := ex.hash()
	txM["txid"] = hex.EncodeToString(h[:])
	txM["type"] = type_exchange
	txM["tx0"] = make(map[string]interface{})
	txM["tx1"] = make(map[string]interface{})

	txM["tx0"].(map[string]interface{})["from"] = ex[0].from
	txM["tx0"].(map[string]interface{})["to"] = ex[0].to
	txM["tx0"].(map[string]interface{})["amount"] = ex[0].amount
	txM["tx0"].(map[string]interface{})["assetId"] = ex[0].assetId
	txM["tx0"].(map[string]interface{})["nonce"] = ex[0].nonce
	txM["tx0"].(map[string]interface{})["fee"] = ex[0].fee
	txM["tx0"].(map[string]interface{})["signature"] = hex.EncodeToString(ex[0].signer.signature[:])

	txM["tx1"].(map[string]interface{})["from"] = ex[1].from
	txM["tx1"].(map[string]interface{})["to"] = ex[1].to
	txM["tx1"].(map[string]interface{})["amount"] = ex[1].amount
	txM["tx1"].(map[string]interface{})["assetId"] = ex[1].assetId
	txM["tx1"].(map[string]interface{})["nonce"] = ex[1].nonce
	txM["tx1"].(map[string]interface{})["fee"] = ex[1].fee
	txM["tx1"].(map[string]interface{})["signature"] = hex.EncodeToString(ex[1].signer.signature[:])

	return txM
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
		"%s", ex.hash(), ex[0].from, ex[0].to, ex[0].amount, ex[0].assetId, ex[0].nonce, ex[0].signer, ex[1].from, ex[1].to, ex[1].amount, ex[1].assetId, ex[1].nonce, ex[1].signer)
}
