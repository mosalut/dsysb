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
	exchange_length = 559
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
		if !validateAddress(transfer.from) {
			return errors.New("`from`: invalid address")
		}

		if !validateAddress(transfer.to) {
			return errors.New("`to`: invalid address")
		}


		if !validateAddress(transfer.hier) {
			return errors.New("`hier`: invalid address")
		}

		if transfer.from == transfer.to {
			return errors.New("Exchange to self is not allowed")
		}

		accountFrom, ok := state.accounts[transfer.from]
		if !ok {
			return errors.New("Transfer from address is not in the state.accounts")
		}

		taskId, ok := state.tasks.isLockedAddress(transfer.from)
		if ok {
			return errors.New("`from`: is locked by task: " + taskId)
		}

		s := hex.EncodeToString(transfer.signer.signature[:])
		if s == "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" {
			return errors.New("Unsigned transaction")
		}

		assetId := fmt.Sprintf("%064x", transfer.assetId)

		fee := transfer.fee()

		if assetId == dsysbId {
			if accountFrom.balance < transfer.amount + fee {
				return errors.New("Not enough DSBs: amount + fee")
			}
		} else {
			asset, ok := state.assets[assetId]
			if !ok {
				return errors.New("No this asset:" + assetId)
			}

			if transfer.bytePrice < asset.price {
				return errors.New(fmt.Sprintf("The byte price should >= asset's create price: %d", asset.price))
			}

			if accountFrom.balance < fee {
				return errors.New("Not enough DSBs: fee")
			}

			balance, ok := accountFrom.assets[assetId]
			if !ok {
				return errors.New("No this asset:" + assetId + " in this account:" + transfer.from)
			}

			if balance < transfer.amount {
				return errors.New("Not enough asset tokens")
			}
		}

		account, ok := state.accounts[transfer.from]
		if !ok {
			return errors.New("No this account:" + transfer.from)
		}

		var nonce uint32
		nonce = account.nonce
		if transfer.nonce - nonce != 1 {
			return errNonceExpired
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

		accountH, ok := state.accounts[transfer.hier]
		if !ok {
			state.accounts[transfer.hier] = &account_T{}
			accountH = state.accounts[transfer.hier]
			accountH.assets = make(map[string]uint64)
		}

		assetId := fmt.Sprintf("%064x", transfer.assetId)
		fee := transfer.fee()

		if assetId == dsysbId {
			if accountFrom.balance < transfer.amount + fee {
				return errors.New("Not enough DSBs: amount + fee")
			}

			accountFrom.balance, accountTo.balance = accountFrom.balance - transfer.amount, accountTo.balance + transfer.amount
		} else {
			if accountFrom.balance < fee {
				return errors.New("not enough DSBs: fee in " + transfer.from)
			}

			balance, ok := accountFrom.assets[assetId]
			if !ok {
				return errors.New("No this asset:" + assetId + " in this account:" + transfer.from)
			}

			if balance < transfer.amount {
				return errors.New("Not enough asset tokens")
			}

			_, ok = accountTo.assets[assetId]
			if !ok {
				accountTo.assets[assetId] = 0
			}
			accountFrom.assets[assetId], accountTo.assets[assetId] = accountFrom.assets[assetId] - transfer.amount, accountTo.assets[assetId] + transfer.amount

			if accountFrom.assets[assetId] == 0 {
				delete(accountFrom.assets, assetId)
			}
		}

		accountFrom.balance -= transfer.fee()
		coinbase.amount += transfer.fee()
		accountFrom.nonce = transfer.nonce
	}

	for _, transfer := range ex {
		accountFrom, _ := state.accounts[transfer.from]
		accountH, _ := state.accounts[transfer.hier]

		if transfer.from == transfer.hier {
			continue
		}

		accountH.balance += accountFrom.balance
		for aid, balance := range accountFrom.assets {
			accountH.assets[aid] += balance
		}

		delete(state.accounts, transfer.from)
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
	txM["tx0"].(map[string]interface{})["hier"] = ex[0].hier
	txM["tx0"].(map[string]interface{})["amount"] = ex[0].amount
	txM["tx0"].(map[string]interface{})["assetId"] = ex[0].assetId
	txM["tx0"].(map[string]interface{})["nonce"] = ex[0].nonce
	txM["tx0"].(map[string]interface{})["fee"] = ex[0].fee
	txM["tx0"].(map[string]interface{})["signature"] = hex.EncodeToString(ex[0].signer.signature[:])

	txM["tx1"].(map[string]interface{})["from"] = ex[1].from
	txM["tx1"].(map[string]interface{})["to"] = ex[1].to
	txM["tx1"].(map[string]interface{})["hier"] = ex[1].hier
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
		"\thier: %s\n" +
		"\tamount: %d\n" +
		"\tasset id: %064x\n" +
		"\tnonce: %d\n" +
		"%s",
		"\tfrom: %s\n" +
		"\tto: %s\n" +
		"\thier: %s\n" +
		"\tamount: %d\n" +
		"\tasset id: %064x\n" +
		"\tnonce: %d\n" +
		"%s", ex.hash(), ex[0].from, ex[0].to, ex[0].hier, ex[0].amount, ex[0].assetId, ex[0].nonce, ex[0].signer, ex[1].from, ex[1].hier, ex[1].to, ex[1].amount, ex[1].assetId, ex[1].nonce, ex[1].signer)
}
