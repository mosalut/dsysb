// dsysb

package main

import (
	"math/big"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	extension_length = 208
	extension_aot_position = 1
	extension_from_position = 2
	extension_nId_position = 36
	extension_blocks_position = 68
	extension_nonce_position = 72
	extension_bytePrice_position = 76
	extension_signer_position = 80
)

type extension_T struct {
	aot byte // asset or task
	from string
	nId [32]byte
	blocks uint32
	nonce uint32
	bytePrice uint32
	signer *signer_T
}

func (et *extension_T) hash() [32]byte {
	bs := make([]byte, extension_signer_position, extension_signer_position)
	bs[0] = type_extension
	bs[extension_aot_position] = et.aot
	copy(bs[extension_from_position:extension_nId_position], et.from)
	copy(bs[extension_nId_position:extension_blocks_position], et.nId[:])
	binary.LittleEndian.PutUint32(bs[extension_blocks_position:extension_nonce_position], et.blocks)
	binary.LittleEndian.PutUint32(bs[extension_nonce_position:extension_bytePrice_position], et.nonce)
	binary.LittleEndian.PutUint32(bs[extension_bytePrice_position:], et.bytePrice)

	return sha256.Sum256(bs)
}

func (et *extension_T) encode() []byte {
	bs := make([]byte, extension_length, extension_length)
	bs[0] = type_extension
	bs[extension_aot_position] = et.aot
	copy(bs[extension_from_position:extension_nId_position], et.from)
	copy(bs[extension_nId_position:extension_blocks_position], et.nId[:])
	binary.LittleEndian.PutUint32(bs[extension_blocks_position:extension_nonce_position], et.blocks)
	binary.LittleEndian.PutUint32(bs[extension_nonce_position:extension_bytePrice_position], et.nonce)
	binary.LittleEndian.PutUint32(bs[extension_bytePrice_position:extension_signer_position], et.bytePrice)
	copy(bs[extension_signer_position:], et.signer.encode())

	return bs
}

func (et *extension_T) encodeForPool() []byte {
	length := extension_length + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], extension_length)
	copy(bs[2:], et.encode())

	return bs
}

func decodeExtension(bs []byte) *extension_T {
	et := &extension_T{}

	et.aot = bs[extension_aot_position]
	et.from = string(bs[extension_from_position:extension_nId_position])
	et.nId = [32]byte(bs[extension_nId_position:extension_blocks_position])
	et.blocks = binary.LittleEndian.Uint32(bs[extension_blocks_position:extension_nonce_position])
	et.nonce = binary.LittleEndian.Uint32(bs[extension_nonce_position:extension_bytePrice_position])
	et.bytePrice = binary.LittleEndian.Uint32(bs[extension_bytePrice_position:extension_signer_position])
	et.signer = decodeSigner(bs[extension_signer_position:])

	return et
}


func (et *extension_T) length() int {
	return extension_length
}

func (et *extension_T) fee() uint64 {
	return extension_length * uint64(et.bytePrice)
}

func (et *extension_T) verifySign() bool {
	publicKey := ecdsa.PublicKey{elliptic.P256(), et.signer.x, et.signer.y}
	txid := et.hash()
	return ecdsa.Verify(&publicKey, txid[:], big.NewInt(0).SetBytes(et.signer.signature[:32]), big.NewInt(0).SetBytes(et.signer.signature[32:]))
}

func (et *extension_T) String() string {
	return fmt.Sprintf(
		"\ttxid:\t%064x\n" +
		"\ttype: extension\n" +
		"\taot: %d\n" +
		"\tfrom: %s\n" +
		"\tnid: %064x\n" +
		"\tnonce: %d\n" +
		"\tbyte price: %d\n" +
		"\tfee: %d\n" +
		"%s", et.hash(), et.aot, et.from, et.nId, et.nonce, et.bytePrice, et.fee(), et.signer)
}

func (et *extension_T) validate(head *blockHead_T, fromP2p bool) error {
	if !validateAddress(et.from) {
		return errors.New("`from`: invalid address")
	}

	/*
	if et.blocks < 10000 {
		return errors.New("`blocks` must >= 10000")
	}
	*/

	if et.bytePrice == 0 {
		return errors.New("Disallow zero byte price")
	}

	// replay attack
	txIdH := et.hash()
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

	id := hex.EncodeToString(et.nId[:])

	switch et.aot {
	case 0:
		_, ok := state.assets[id]
		if !ok {
			return errors.New("There is not this asset")
		}
	case 1:
		var b bool
		for _, t := range state.tasks {
			tId := t.hash()
			if hex.EncodeToString(tId[:]) == id {
				b = true
				break
			}
		}
		if !b {
			return errors.New("There is not this task")
		}
	}

	account, ok := state.accounts[et.from]
	if !ok {
		return errors.New("There's not the account id")
	}

	nonce = account.nonce
	fmt.Println("nonce:", et.nonce, nonce)
	if et.nonce - nonce != 1 {
		return errNonceExpired
	}

	ok = et.verifySign()
	if !ok {
		return errors.New("Invalid signature")
	}

	return nil
}

func (et *extension_T) count(state *state_T, coinbase *coinbase_T, index int) error {
	account, ok := state.accounts[et.from]
	if !ok {
		return errors.New("DT address is empty address")
	}

	id := hex.EncodeToString(et.nId[:])

	switch et.aot {
	case 0:
		asset, ok := state.assets[id]
		if !ok {
			return errors.New("There is not this asset")
		}
		holdAmount := uint64(asset.price) * uint64(et.blocks)
		totalSpend := holdAmount + et.fee()
		if account.balance < totalSpend {
			return errors.New("not enough DSBs")
		}
		account.balance -= totalSpend
		asset.remain += et.blocks
	case 1:
		var task *task_T
		for k, t := range state.tasks {
			tId := t.hash()
			if hex.EncodeToString(tId[:]) == id {
				task = state.tasks[k]
			}
		}
		holdAmount := uint64(task.price) * uint64(et.blocks)
		totalSpend := holdAmount + et.fee()
		if account.balance < totalSpend {
			return errors.New("not enough DSBs")
		}
		account.balance -= totalSpend
		task.remain += et.blocks
	}

	account.nonce = et.nonce

	return nil
}

func (et *extension_T) Map() map[string]interface{} {
	txM := make(map[string]interface{})
	h := et.hash()
	txM["txid"] = hex.EncodeToString(h[:])
	txM["type"] = type_extension
	txM["aot"] = et.aot
	txM["from"] = et.from
	txM["nId"] = hex.EncodeToString(et.nId[:])
	txM["nonce"] = et.nonce
	txM["byte price"] = et.bytePrice
	txM["fee"] = et.fee()
	txM["signature"] = hex.EncodeToString(et.signer.signature[:])

	return txM
}

func (et *extension_T) getBytePrice() uint32 {
	return et.bytePrice
}
