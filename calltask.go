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

type callTask_T struct {
	taskId [32]byte
	from string
	params []byte
	nonce uint32
	fee uint64
	signer *signer_T
}

func (tx *callTask_T) hash() [32]byte {
	pLength := len(tx.params)
	paramsEnd := 66 + pLength // 66 = 32 + 34
	nonceEnd := paramsEnd + 4
	feeEnd := nonceEnd + 8
	length := feeEnd + 128
	bs := make([]byte, length, length)
	copy(bs[:32], tx.taskId[:])
	copy(bs[32:66], []byte(tx.from))
	copy(bs[66:paramsEnd], tx.params)
	binary.LittleEndian.PutUint32(bs[paramsEnd:nonceEnd], tx.nonce)
	binary.LittleEndian.PutUint64(bs[nonceEnd:feeEnd], tx.fee)

	return sha256.Sum256(bs)
}

func (tx *callTask_T) encode() []byte {
	pLength := len(tx.params)
	paramsEnd := 66 + pLength // 66 = 32 + 34
	nonceEnd := paramsEnd + 4
	feeEnd := nonceEnd + 8
	length := feeEnd + 128
	bs := make([]byte, length, length)
	copy(bs[:32], tx.taskId[:])
	copy(bs[32:66], []byte(tx.from))
	copy(bs[66:paramsEnd], tx.params)
	binary.LittleEndian.PutUint32(bs[paramsEnd:nonceEnd], tx.nonce)
	binary.LittleEndian.PutUint64(bs[nonceEnd:feeEnd], tx.fee)
	copy(bs[feeEnd:], tx.signer.encode())

	return bs
}

func (tx *callTask_T) encodeForPool() []byte {
	// 206 = 32 + 34 + 4 + 8 + 128
	length0 := 206 + len(tx.params)
	length := length0 + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], uint16(length0))
	copy(bs[2:], tx.encode())

	return bs
}

func decodeCallTask(bs []byte) *callTask_T {
	tx := &callTask_T{}
	tx.taskId = [32]byte(bs[:32])
	tx.from = string(bs[32:66])
	paramsEnd := len(bs) - 140 // 140 = 4 + 8 + 128
	tx.params = bs[66:paramsEnd]
	nonceEnd := paramsEnd + 4
	tx.nonce = binary.LittleEndian.Uint32(bs[paramsEnd:nonceEnd])
	feeEnd := nonceEnd + 8
	tx.fee = binary.LittleEndian.Uint64(bs[nonceEnd:feeEnd])
	tx.signer = decodeSigner(bs[feeEnd:])

	return tx
}

func (ct *callTask_T) validate(fromP2p bool) error {
	txIdsMutex.Lock()
	defer txIdsMutex.Unlock()

	if len(ct.params) > 65330 {
		return errors.New("warning: params's length is too long")
	}

	if ct.fee == 0 {
		fmt.Println("warning: got zero fee")
	}

	if !validateAddress(ct.from) {
		return errors.New("`from`: invalid address")
	}

	s := hex.EncodeToString(ct.signer.signature[:])
	if s == "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" {
		return errors.New("Unsigned transaction")
	}

	// replay attack
	txIdH := ct.hash()
	txId := hex.EncodeToString(txIdH[:])
	for _, id := range txIds {
		if txId == id {
			if fromP2p {
				deleteFromTransactionPool(txId)
				return nil
			}
			return errors.New("Replay attack: txid:" + txId)
		}
	}

	var nonce uint32
	state, err := getState()
	if err != nil {
		return err
	}

	account, ok := state.accounts[ct.from]
	if ok {
		nonce = account.nonce
	}

	fmt.Println("nonce:", ct.nonce, nonce)
	if ct.nonce - nonce != 1 {
		return errors.New("The nonces are not match")
	}

	ok = ct.verifySign()
	if !ok {
		return errors.New("Invalid signature")
	}

	txIds = append(txIds, txId)

	return nil
}

func (ct *callTask_T) verifySign() bool {
	publicKey := ecdsa.PublicKey{elliptic.P256(), ct.signer.x, ct.signer.y}
	txid := ct.hash()
	return ecdsa.Verify(&publicKey, txid[:], big.NewInt(0).SetBytes(ct.signer.signature[:32]), big.NewInt(0).SetBytes(ct.signer.signature[32:]))
}

func (ct *callTask_T) count(state *state_T, coinbase *coinbase_T, index int) error {
	var task *task_T
	for k, t := range state.tasks {
		tId := t.hash()
		if hex.EncodeToString(tId[:]) == hex.EncodeToString(ct.taskId[:]) {
			task = state.tasks[k]
		}
	}

	if task == nil {
		return errors.New("The task of CT is not found")
	}

	err := task.excute(state)
	if err != nil {
		return err
	}

	account, ok := state.accounts[ct.from]
	if !ok {
		return errors.New("CT address is empty address")
	}

	if account.balance < ct.fee {
		return errors.New("not enough minerals")
	}
	account.balance -= ct.fee
	coinbase.amount += ct.fee
	account.nonce = ct.nonce

	return nil
}

func (tx *callTask_T) Map() map[string]interface{} {
	txM := make(map[string]interface{})
	h := tx.hash()
	txM["txid"] = hex.EncodeToString(h[:])
	txM["type"] = type_call
	txM["taskId"] = hex.EncodeToString(tx.taskId[:])
	txM["from"] = tx.from
	txM["params"] = hex.EncodeToString(tx.params[:])
	txM["nonce"] = tx.nonce
	txM["fee"] = tx.fee
	txM["signature"] = hex.EncodeToString(tx.signer.signature[:])

	return txM
}

func (tx *callTask_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
			"\ttype:\tcall\n" +
			"\ttask id: %x\n" +
			"\tfrom: %s\n" +
			"\tparams: %v\n" +
			"\tnonce: %d\n" +
			"\tfee: %d\n" +
			"%s", tx.hash(), tx.taskId, tx.from, tx.params, tx.nonce, tx.fee, tx.signer)
}

func isCall(bs []byte) bool {
	length := len(bs)
	if length < 206 || length > 65536 {
		return false
	}

	return validateAddress(string(bs[32:66]))
}
