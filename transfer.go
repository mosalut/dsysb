// dsysb

package main

import (
	"math/big"
	"crypto/sha256"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	transfer_length = 248
	transfer_to_position = 34
	transfer_amount_position = 68
	transfer_asset_id_position = 76
	transfer_nonce_position = 108
	transfer_fee_position = 112
	transfer_signer_position = 120
)

type transfer_T struct {
	from string
	to string
	amount uint64
	assetId [32]byte
	nonce uint32
	fee uint64
	signer *signer_T
}

func (transfer *transfer_T) hash() [32]byte {
	bs := transfer.encodeWithoutSigner()

	return sha256.Sum256(bs)
}

func (transfer *transfer_T) encode() []byte {
	bs := make([]byte, transfer_length, transfer_length)
	copy(bs[:transfer_to_position], []byte(transfer.from))
	copy(bs[transfer_to_position:transfer_amount_position], []byte(transfer.to))
	binary.LittleEndian.PutUint64(bs[transfer_amount_position:transfer_asset_id_position],transfer.amount)
	copy(bs[transfer_asset_id_position:transfer_nonce_position], transfer.assetId[:])
	binary.LittleEndian.PutUint32(bs[transfer_nonce_position:transfer_fee_position], transfer.nonce)
	binary.LittleEndian.PutUint64(bs[transfer_fee_position:transfer_signer_position], transfer.fee)
	copy(bs[transfer_signer_position:], transfer.signer.encode())

	return bs
}

func (transfer *transfer_T) encodeForPool() []byte {
	length := transfer_length + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], transfer_length)
	copy(bs[2:], transfer.encode())

	return bs
}

func decodeTransfer(bs []byte) *transfer_T {
	transfer := &transfer_T{}
	transfer.from = string(bs[:transfer_to_position])
	transfer.to = string(bs[transfer_to_position:transfer_amount_position])
	transfer.amount = binary.LittleEndian.Uint64(bs[transfer_amount_position:transfer_asset_id_position])
	transfer.assetId = [32]byte(bs[transfer_asset_id_position:transfer_nonce_position])
	transfer.nonce = binary.LittleEndian.Uint32(bs[transfer_nonce_position:transfer_fee_position])
	transfer.fee = binary.LittleEndian.Uint64(bs[transfer_fee_position:transfer_signer_position])
	transfer.signer = decodeSigner(bs[transfer_signer_position:])

	return transfer
}

func (transfer *transfer_T) encodeWithoutSigner() []byte {
	bs := make([]byte, transfer_signer_position, transfer_signer_position)
	copy(bs[:transfer_to_position], []byte(transfer.from))
	copy(bs[transfer_to_position:transfer_amount_position], []byte(transfer.to))
	binary.LittleEndian.PutUint64(bs[transfer_amount_position:transfer_asset_id_position],transfer.amount)
	copy(bs[transfer_asset_id_position:transfer_nonce_position], transfer.assetId[:])
	binary.LittleEndian.PutUint32(bs[transfer_nonce_position:transfer_fee_position], transfer.nonce)
	binary.LittleEndian.PutUint64(bs[transfer_fee_position:transfer_signer_position], transfer.fee)

	return bs
}

func (transfer *transfer_T) validate(fromP2p bool) error {
	txIdsMutex.Lock()
	defer txIdsMutex.Unlock()

	if transfer.from == transfer.to {
		return errors.New("Transfer to self is not allowed")
	}

	s := hex.EncodeToString(transfer.signer.signature[:])
	if s == "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" {
		return errors.New("Unsigned transaction")
	}

	poolMutex.Lock()
	defer poolMutex.Unlock()

	// replay attack
	txIdH := transfer.hash()
	txId := hex.EncodeToString(txIdH[:])
	for _, id := range txIds {
		if txId == id {
			if fromP2p {
				deleteFromTransactionPool(txId)
				return nil
			}
			return errors.New("Replay attack: txid: " + txId)
		}
	}

	state, err := getState()
	if err != nil {
		return err
	}

	assetId := hex.EncodeToString(transfer.assetId[:])

	if assetId != dsysbId {
		_, ok := state.assets[assetId]
		if !ok {
			print(log_error, "There's not the asset id: " + assetId)
			return errors.New("There's not the asset id: " + assetId)
		}
	}

	account, ok := state.accounts[transfer.from]
	if !ok {
		return errors.New("There's not the account id")
	}

	nonce := account.nonce
	if transfer.nonce - nonce != 1 {
		return errors.New("The nonces are not match")
	}

	ok = transfer.verifySign()
	if !ok {
		return errors.New("Invalid signature")
	}

	txIds = append(txIds, txId)

	return nil
}

func (transfer *transfer_T) verifySign() bool {
	publicKey := ecdsa.PublicKey{elliptic.P256(), transfer.signer.x, transfer.signer.y}
	txid := transfer.hash()
	return ecdsa.Verify(&publicKey, txid[:], big.NewInt(0).SetBytes(transfer.signer.signature[:32]), big.NewInt(0).SetBytes(transfer.signer.signature[32:]))
}

func (transfer *transfer_T) count(state *state_T, coinbase *coinbase_T, index int) error {
	accountFrom, ok := state.accounts[transfer.from]
	if !ok {
		return errors.New("The address of transfer from is empty")
	}

	accountTo, ok := state.accounts[transfer.to]
	if !ok {
		state.accounts[transfer.to] = &account_T{}
		accountTo = state.accounts[transfer.to]
		accountTo.assets = make(map[string]uint64)
	}

	id := fmt.Sprintf("%064x", transfer.assetId)

	if id == dsysbId {
		if accountFrom.balance < transfer.amount + transfer.fee {
			return errors.New("not enough DSBs: amount + fee")
		}

		accountFrom.balance, accountTo.balance = accountFrom.balance - transfer.amount, accountTo.balance + transfer.amount
	} else {
		if accountFrom.balance < transfer.fee {
			return errors.New("not enough DSBs: fee")
		}
		balance, ok := accountFrom.assets[id]
		if !ok {
			return errors.New("There is not this asset")
		}

		if balance < transfer.amount {
			return errors.New("not enough asset token")
		}

		_, ok = accountTo.assets[id]
		if !ok {
			accountTo.assets[id] = 0
		}
		accountFrom.assets[id], accountTo.assets[id] = accountFrom.assets[id] - transfer.amount, accountTo.assets[id] + transfer.amount
	}

	accountFrom.balance -= transfer.fee
//	state.accounts[*address].balance += transfer.fee
	coinbase.amount += transfer.fee
	accountFrom.nonce = transfer.nonce

	return nil
}

func (transfer *transfer_T) String() string {
	return fmt.Sprintf(
		"\tfrom: %s\n" +
		"\tto: %s\n" +
		"\tamount: %d\n" +
		"\tasset id: %064x\n" +
		"\tnonce: %d\n" +
		"\tfee: %d\n" +
		"%s",
		transfer.from, transfer.to, transfer.amount, transfer.assetId, transfer.nonce, transfer.fee, transfer.signer)
}
