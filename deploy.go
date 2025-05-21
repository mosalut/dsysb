// dsysb

package main

import (
	"math/big"
	"crypto/sha256"
	"crypto/elliptic"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

type deployTask_T struct {
	instructs []uint8
	vData []byte
	from string
	nonce uint32
	fee uint64
	signer *signer_T
}

func (tx *deployTask_T) hash() [32]byte {
	// 50 = 2 + 2 + 34 + 4 + 8
	length := 50 + len(tx.instructs) + len(tx.vData)
	bs := make([]byte, length, length)
	instructsLength := len(tx.instructs)
	binary.LittleEndian.PutUint16(bs[:2], uint16(instructsLength))
	vDataLengthPosition := 2 + instructsLength
	copy(bs[2:vDataLengthPosition], tx.instructs)
	vDataLength := len(tx.vData)
	vDataPosition := vDataLengthPosition + 2
	binary.LittleEndian.PutUint16(bs[vDataLengthPosition:vDataPosition], uint16(vDataLength))
	fromPosition := vDataPosition + vDataLength
	copy(bs[vDataPosition:fromPosition], tx.vData)
	noncePosition := fromPosition + 34
	copy(bs[fromPosition:noncePosition], tx.from)
	feePosition := noncePosition + 4
	binary.LittleEndian.PutUint32(bs[noncePosition:feePosition], uint32(tx.nonce))
	signerPosition := feePosition + 8
	binary.LittleEndian.PutUint64(bs[feePosition:signerPosition], uint64(tx.fee))

	return sha256.Sum256(bs)
}

func (tx *deployTask_T) encode() []byte {
	// 178 = 2 + 2 + 34 + 4 + 8 + 128
	length := 178 + len(tx.instructs) + len(tx.vData)
	bs := make([]byte, length, length)
	instructsLength := len(tx.instructs)
	binary.LittleEndian.PutUint16(bs[:2], uint16(instructsLength))
	vDataLengthPosition := 2 + instructsLength
	copy(bs[2:vDataLengthPosition], tx.instructs)
	vDataLength := len(tx.vData)
	vDataPosition := vDataLengthPosition + 2
	binary.LittleEndian.PutUint16(bs[vDataLengthPosition:vDataPosition], uint16(vDataLength))
	fromPosition := vDataPosition + vDataLength
	copy(bs[vDataPosition:fromPosition], tx.vData)
	noncePosition := fromPosition + 34
	copy(bs[fromPosition:noncePosition], tx.from)
	feePosition := noncePosition + 4
	binary.LittleEndian.PutUint32(bs[noncePosition:feePosition], uint32(tx.nonce))
	signerPosition := feePosition + 8
	binary.LittleEndian.PutUint64(bs[feePosition:signerPosition], uint64(tx.fee))
	copy(bs[signerPosition:], tx.signer.encode())

	return bs
}

func (tx *deployTask_T) encodeForPool() []byte {
	length0 := 178 + len(tx.instructs) + len(tx.vData)
	length := length0 + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], uint16(length0))
	copy(bs[2:], tx.encode())

	return bs
}

func decodeDeployTask(bs []byte) *deployTask_T {
	tx := &deployTask_T{}
	instructsLength := binary.LittleEndian.Uint16(bs[:2])
	vDataLengthPosition := 2 + instructsLength
	tx.instructs = bs[2:vDataLengthPosition]
	vDataPosition := vDataLengthPosition + 2
	vDataLength := binary.LittleEndian.Uint16(bs[vDataLengthPosition:vDataPosition])
	fromPosition := vDataPosition + vDataLength
	tx.vData = bs[vDataPosition:fromPosition]
	noncePosition := fromPosition + 34
	tx.from = string(bs[fromPosition:noncePosition])
	feePosition := noncePosition + 4
	tx.nonce = binary.LittleEndian.Uint32(bs[noncePosition:feePosition])
	signerPosition := feePosition + 8
	tx.fee = binary.LittleEndian.Uint64(bs[feePosition:signerPosition])
	tx.signer = decodeSigner(bs[signerPosition:])

	return tx
}

func (tx *deployTask_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
			"\ttype:\tdeploy\n" +
			"\tinstructs: %v\n" +
			"\tvData: %v\n" +
			"\tfrom: %s\n" +
			"\tnonce: %d\n" +
			"\tfee: %d\n" +
			"%s", tx.hash(), tx.instructs, tx.vData, tx.from, tx.nonce, tx.fee, tx.signer)
}

func isDeploy(bs []byte) bool {
	if(len(bs) < 2) {
		return false
	}

	instructsLength := binary.LittleEndian.Uint16(bs[:2])
	vDataLengthPosition := 2 + instructsLength

	if(uint16(len(bs)) < vDataLengthPosition) {
		return false
	}

	vDataPosition := vDataLengthPosition + 2
	if(uint16(len(bs)) < vDataPosition) {
		return false
	}

	vDataLength := binary.LittleEndian.Uint16(bs[vDataLengthPosition:vDataPosition])
	fromPosition := vDataPosition + vDataLength

	if(uint16(len(bs)) != fromPosition + 174) {
		return false
	}

	return true
}

func (dt *deployTask_T) validate(fromP2p bool) error {
	txIdsMutex.Lock()
	defer txIdsMutex.Unlock()

	if len(dt.instructs) + len(dt.vData) > 65358 {
		return errors.New("warning: instructs' and vdata's length is too long")
	}

	if dt.fee == 0 {
		fmt.Println("warning: got zero fee")
	}

	if !validateAddress(dt.from) {
		return errors.New("`from`: invalid address")
	}

	s := hex.EncodeToString(dt.signer.signature[:])
	if s == "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" {
		return errors.New("Unsigned transaction")
	}

	// replay attack
	txIdH := dt.hash()
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

	account, ok := state.accounts[dt.from]
	if ok {
		nonce = account.nonce
	}

	fmt.Println("nonce:", dt.nonce, nonce)
	if dt.nonce - nonce != 1 {
		return errors.New("The nonces are not match")
	}

	ok = dt.verifySign()
	if !ok {
		return errors.New("Invalid signature")
	}

	txIds = append(txIds, txId)

	return nil
}

func (dt *deployTask_T) verifySign() bool {
	publicKey := ecdsa.PublicKey{elliptic.P256(), dt.signer.x, dt.signer.y}
	txid := dt.hash()
	return ecdsa.Verify(&publicKey, txid[:], big.NewInt(0).SetBytes(dt.signer.signature[:32]), big.NewInt(0).SetBytes(dt.signer.signature[32:]))
}

func (dt *deployTask_T) count(state *state_T, coinbase *coinbase_T, index int) error {
	task := &task_T {
		dt.from,
		dt.instructs,
		dt.vData,
	}

	taskIdB := task.hash()
	taskId := fmt.Sprintf("%064x", taskIdB)
	fmt.Println("taskId:", taskId)
	for _, task := range state.tasks {
		h := task.hash()
		if hex.EncodeToString(h[:]) == taskId {
			return errors.New("task is already in")
		}
	}

	account, ok := state.accounts[task.address]
	if !ok {
		return errors.New("DT address is empty address")
	}

	if account.balance < dt.fee {
		return errors.New("not enough minerals")
	}

	state.tasks = append(state.tasks, task)

	account.balance -= dt.fee
	coinbase.amount += dt.fee
	account.nonce = dt.nonce

	return nil
}
