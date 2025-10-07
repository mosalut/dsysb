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
	from string
	hier string
	taskId [32]byte
	params []byte
	nonce uint32
	bytePrice uint32
	signer *signer_T
}

func (tx *callTask_T) hash() [32]byte {
	pLength := len(tx.params)
	paramsEnd := 101 + pLength // 101 = 1 + 32 + 34 + 34
	nonceEnd := paramsEnd + 4
	bytePriceEnd := nonceEnd + 4
	length := bytePriceEnd + 128
	bs := make([]byte, length, length)
	bs[0] = type_call
	copy(bs[1:35], []byte(tx.from))
	copy(bs[35:69], []byte(tx.hier))
	copy(bs[69:101], tx.taskId[:])
	copy(bs[101:paramsEnd], tx.params)
	binary.LittleEndian.PutUint32(bs[paramsEnd:nonceEnd], tx.nonce)
	binary.LittleEndian.PutUint32(bs[nonceEnd:bytePriceEnd], tx.bytePrice)

	return sha256.Sum256(bs)
}

func (tx *callTask_T) encode() []byte {
	pLength := len(tx.params)
	paramsEnd := 101 + pLength // 101 = 1 + 32 + 34 + 34
	nonceEnd := paramsEnd + 4
	bytePriceEnd := nonceEnd + 4
	length := bytePriceEnd + 128
	bs := make([]byte, length, length)
	bs[0] = type_call
	copy(bs[1:35], []byte(tx.from))
	copy(bs[35:69], []byte(tx.hier))
	copy(bs[69:101], tx.taskId[:])
	copy(bs[101:paramsEnd], tx.params)
	binary.LittleEndian.PutUint32(bs[paramsEnd:nonceEnd], tx.nonce)
	binary.LittleEndian.PutUint32(bs[nonceEnd:bytePriceEnd], tx.bytePrice)
	copy(bs[bytePriceEnd:], tx.signer.encode())

	return bs
}

func (tx *callTask_T) encodeForPool() []byte {
	// 237 = 1 + 32 + 34 + 34 + 4 + 4 + 128
	length0 := 237 + len(tx.params)
	length := length0 + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], uint16(length0))
	copy(bs[2:], tx.encode())

	return bs
}

func decodeCallTask(bs []byte) *callTask_T {
	tx := &callTask_T{}
	tx.from = string(bs[1:35])
	tx.hier = string(bs[35:69])
	tx.taskId = [32]byte(bs[69:101])
	paramsEnd := len(bs) - 136 // 136 = 4 + 4 + 128
	tx.params = bs[101:paramsEnd]
	nonceEnd := paramsEnd + 4
	tx.nonce = binary.LittleEndian.Uint32(bs[paramsEnd:nonceEnd])
	bytePriceEnd := nonceEnd + 4
	tx.bytePrice = binary.LittleEndian.Uint32(bs[nonceEnd:bytePriceEnd])
	tx.signer = decodeSigner(bs[bytePriceEnd:])

	return tx
}

func (tx *callTask_T) length() int {
	return len(tx.params) + 237
}

func (tx *callTask_T) fee() uint64 {
	return uint64(tx.length()) * uint64(tx.bytePrice)
}

func (ct *callTask_T) validate(head *blockHead_T, fromP2p bool) error {
	if len(ct.params) > 65333 {
		return errors.New("Params's length is too long")
	}

	if ct.fee() == 0 {
		return errors.New("Disallow zero byte price")
	}

	if !validateAddress(ct.from) {
		return errors.New("`from`: invalid address")
	}

	if !validateAddress(ct.hier) {
		return errors.New("`hier`: invalid address")
	}

	s := hex.EncodeToString(ct.signer.signature[:])
	if s == "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" {
		return errors.New("Unsigned transaction")
	}

	// replay attack
	txIdH := ct.hash()
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

	state, err := getState()
	if err != nil {
		return err
	}

	var ok bool
	for _, task := range state.tasks {
		h := task.hash()
		if hex.EncodeToString(h[:]) == hex.EncodeToString(ct.taskId[:]) {
			if ct.bytePrice < task.price {
				return errors.New(fmt.Sprintf("The byte price should >= task's deploy price: %d", task.price))
			}
			err := task.validateCall(state, ct)
			if err != nil {
				return err
			}
			ok = true
			break
		}
	}

	if !ok {
		return errors.New("The task of CT is not found")
	}

	account, ok := state.accounts[ct.from]
	if !ok {
		return errors.New("CT address is empty address")
	}
	nonce := account.nonce

	if account.balance < ct.fee() {
		return errors.New("Not enough minerals")
	}

	fmt.Println("nonce:", ct.nonce, nonce)
	if ct.nonce - nonce != 1 {
		return errNonceExpired
	}

	ok = ct.verifySign()
	if !ok {
		return errors.New("Invalid signature")
	}

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
			break
		}
	}

	if task == nil {
		return errors.New("The task of CT is not found")
	}

	account, ok := state.accounts[ct.from]
	if !ok {
		return errors.New("CT address is empty address")
	}

	accountH, ok := state.accounts[ct.hier]
	if !ok {
		state.accounts[ct.hier] = &account_T{}
		accountH = state.accounts[ct.hier]
		accountH.assets = make(map[string]uint64)
	}

	if account.balance < ct.fee() {
		return errors.New("Not enough minerals")
	}

	account.balance -= ct.fee()
	coinbase.amount += ct.fee()
	account.nonce = ct.nonce

	err := task.excute(state, ct.from, ct.fee(), ct.params)
	if err != nil {
		print(log_warning, "taskId:", hex.EncodeToString(ct.taskId[:]), "caller:", ct.from, "excute error:", err)
		return errors.New("taskId:" + hex.EncodeToString(ct.taskId[:]) + "caller:" + ct.from + " excute error:" + err.Error())
	}

	if ct.from == ct.hier {
		return nil
	}

	accountH.balance += account.balance
	for aid, balance := range account.assets {
		accountH.assets[aid] += balance
	}

	return nil
}

func (tx *callTask_T) getBytePrice() uint32 {
	return tx.bytePrice
}

func (tx *callTask_T) Map() map[string]interface{} {
	txM := make(map[string]interface{})
	h := tx.hash()
	txM["txid"] = hex.EncodeToString(h[:])
	txM["type"] = type_call
	txM["from"] = tx.from
	txM["hier"] = tx.hier
	txM["taskId"] = hex.EncodeToString(tx.taskId[:])
	txM["params"] = hex.EncodeToString(tx.params[:])
	txM["nonce"] = tx.nonce
	txM["byte price"] = tx.bytePrice
	txM["fee"] = tx.fee()
	txM["signature"] = hex.EncodeToString(tx.signer.signature[:])

	return txM
}

func (tx *callTask_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
			"\ttype:\tcall\n" +
			"\tfrom: %s\n" +
			"\thier: %s\n" +
			"\ttask id: %x\n" +
			"\tparams: %v\n" +
			"\tnonce: %d\n" +
			"\tbyte price: %d\n" +
			"\tfee: %d\n" +
			"%s", tx.hash(), tx.from, tx.hier, tx.taskId, tx.params, tx.nonce, tx.bytePrice, tx.fee(), tx.signer)
}
