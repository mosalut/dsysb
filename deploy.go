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

const (
	dtLengthWithoutSignature = 55 // 55 = type:1 + len_ins:2 + len_v:2 + price:4 + blocks:4 + from:34 + nonce:4 + bytePrice:4
	dtLength = 183 // 183 = type:1 + len_ins:2 + len_v:2 + price:4 + blocks:4 + from:34 + nonce:4 + bytePrice:4 + signature:128
)

type deployTask_T struct {
	instructs []uint8
	vData []byte
	price uint32
	blocks uint32
	from string
	nonce uint32
	bytePrice uint32
	signer *signer_T
}

func (tx *deployTask_T) hash() [32]byte {
	length := dtLengthWithoutSignature + len(tx.instructs) + len(tx.vData)
	bs := make([]byte, length, length)
	bs[0] = type_deploy
	instructsLength := len(tx.instructs)
	binary.LittleEndian.PutUint16(bs[1:3], uint16(instructsLength))
	vDataLengthPosition := 3 + instructsLength
	copy(bs[3:vDataLengthPosition], tx.instructs)
	vDataLength := len(tx.vData)
	vDataPosition := vDataLengthPosition + 2
	binary.LittleEndian.PutUint16(bs[vDataLengthPosition:vDataPosition], uint16(vDataLength))
	pricePosition := vDataPosition + vDataLength
	copy(bs[vDataPosition:pricePosition], tx.vData)
	binary.LittleEndian.PutUint32(bs[pricePosition:pricePosition + 4], tx.price)
	binary.LittleEndian.PutUint32(bs[pricePosition + 4:pricePosition + 8], tx.blocks)
	fromPosition := pricePosition + 8
	noncePosition := fromPosition + 34
	copy(bs[fromPosition:noncePosition], tx.from)
	bytePricePosition := noncePosition + 4
	binary.LittleEndian.PutUint32(bs[noncePosition:bytePricePosition], tx.nonce)
	signerPosition := bytePricePosition + 4
	binary.LittleEndian.PutUint32(bs[bytePricePosition:signerPosition], tx.bytePrice)

	return sha256.Sum256(bs)
}

func (tx *deployTask_T) encode() []byte {
	length := dtLength + len(tx.instructs) + len(tx.vData)
	bs := make([]byte, length, length)
	bs[0] = type_deploy
	instructsLength := len(tx.instructs)
	binary.LittleEndian.PutUint16(bs[1:3], uint16(instructsLength))
	vDataLengthPosition := 3 + instructsLength
	copy(bs[3:vDataLengthPosition], tx.instructs)
	vDataLength := len(tx.vData)
	vDataPosition := vDataLengthPosition + 2
	binary.LittleEndian.PutUint16(bs[vDataLengthPosition:vDataPosition], uint16(vDataLength))

	pricePosition := vDataPosition + vDataLength
	copy(bs[vDataPosition:pricePosition], tx.vData)

	binary.LittleEndian.PutUint32(bs[pricePosition:pricePosition + 4], tx.price)
	binary.LittleEndian.PutUint32(bs[pricePosition + 4:pricePosition + 8], tx.blocks)

	fromPosition := pricePosition + 8
	noncePosition := fromPosition + 34
	copy(bs[fromPosition:noncePosition], tx.from)
	bytePricePosition := noncePosition + 4
	binary.LittleEndian.PutUint32(bs[noncePosition:bytePricePosition], tx.nonce)
	signerPosition := bytePricePosition + 4
	binary.LittleEndian.PutUint32(bs[bytePricePosition:signerPosition], tx.bytePrice)
	copy(bs[signerPosition:], tx.signer.encode())

	return bs
}

func (tx *deployTask_T) encodeForPool() []byte {
	length0 := dtLength + len(tx.instructs) + len(tx.vData)
	length := length0 + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], uint16(length0))
	copy(bs[2:], tx.encode())

	return bs
}

func decodeDeployTask(bs []byte) *deployTask_T {
	tx := &deployTask_T{}
	instructsLength := binary.LittleEndian.Uint16(bs[1:3])
	vDataLengthPosition := 3 + instructsLength
	tx.instructs = bs[3:vDataLengthPosition]
	vDataPosition := vDataLengthPosition + 2
	vDataLength := binary.LittleEndian.Uint16(bs[vDataLengthPosition:vDataPosition])
	pricePosition := vDataPosition + vDataLength
	tx.vData = bs[vDataPosition:pricePosition]
	tx.price = binary.LittleEndian.Uint32(bs[pricePosition:pricePosition + 4])
	tx.blocks = binary.LittleEndian.Uint32(bs[pricePosition + 4:pricePosition + 8])
	fromPosition := pricePosition + 8
	noncePosition := fromPosition + 34
	tx.from = string(bs[fromPosition:noncePosition])
	bytePricePosition := noncePosition + 4
	tx.nonce = binary.LittleEndian.Uint32(bs[noncePosition:bytePricePosition])
	signerPosition := bytePricePosition + 4
	tx.bytePrice = binary.LittleEndian.Uint32(bs[bytePricePosition:signerPosition])
	tx.signer = decodeSigner(bs[signerPosition:])

	return tx
}

func (tx *deployTask_T) getBytePrice() uint32 {
	return tx.bytePrice
}

func (tx *deployTask_T) Map() map[string]interface{} {
	txM := make(map[string]interface{})
	h := tx.hash()
	txM["txid"] = hex.EncodeToString(h[:])
	txM["type"] = type_deploy
	txM["instructs"] = hex.EncodeToString(tx.instructs[:])
	txM["vData"] = hex.EncodeToString(tx.vData[:])
	txM["price"] = tx.price
	txM["blocks"] = tx.blocks
	txM["from"] = tx.from
	txM["nonce"] = tx.nonce
	txM["bytePrice"] = tx.bytePrice
	txM["signature"] = hex.EncodeToString(tx.signer.signature[:])

	return txM
}

func (tx *deployTask_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
			"\ttype:\tdeploy\n" +
			"\tinstructs: %v\n" +
			"\tvData: %v\n" +
			"\tprice: %d\n" +
			"\tblocks: %d\n" +
			"\tfrom: %s\n" +
			"\tnonce: %d\n" +
			"\tbyte price: %d\n" +
			"%s", tx.hash(), tx.instructs, tx.vData, tx.price, tx.blocks, tx.from, tx.nonce, tx.bytePrice, tx.signer)
}

func (dt *deployTask_T) validate(head *blockHead_T, fromP2p bool) error {
	if len(dt.instructs) + len(dt.vData) > 65353 {
		return errors.New("Instructs' and vdata's length is too long")
	}

	if dt.bytePrice == 0 {
		return errors.New("Disallow zero byte price")
	}

	if dt.price == 0 {
		return errors.New("Deploy's price must > 0")
	}

	if dt.blocks < 10000 {
		return errors.New("Deploy's blocks must >= 10000")
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

	return nil
}

func (tx *deployTask_T) length() int {
	return dtLength + len(tx.instructs) + len(tx.vData)
}

func (tx *deployTask_T) fee() uint64 {
	return uint64(tx.length()) * uint64(tx.bytePrice)
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
		dt.nonce,
		dt.price,
		dt.blocks,
		dt.blocks,
	}

	taskIdB := task.hash()
	taskId := hex.EncodeToString(taskIdB[:])
	for _, task := range state.tasks {
		h := task.hash()
		if hex.EncodeToString(h[:]) == taskId {
			return errors.New("task is already in")
		}
	}

	account, ok := state.accounts[dt.from]
	if !ok {
		return errors.New("DT address is empty address")
	}

	holdAmount := uint64(dt.price) * uint64(dt.blocks)
	totalSpend := holdAmount + dt.fee()

	if account.balance < totalSpend {
		return errors.New("not enough minerals")
	}

	state.tasks = append(state.tasks, task)

	account.balance -= totalSpend
	coinbase.amount += dt.fee()
	account.nonce = dt.nonce

	return nil
}
