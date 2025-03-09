// dsysb

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

type account_T struct {
	balance uint64
	assets map[string]uint64 // key is an asset id
	nonce uint32
}

func (account *account_T)encode() []byte {
	length := 12 + len(account.assets) * 40
	bs := make([]byte, length, length)
	var start int
	end := 8

	binary.LittleEndian.PutUint64(bs[:end], account.balance)

	for k, asset := range account.assets {
		start = end
		end += 32
		key, err := hex.DecodeString(k)
		if err != nil {
			print(log_error, err)
			return nil
		}
		copy(bs[start:end], key)
		start = end
		end += 8
		binary.LittleEndian.PutUint64(bs[start:end], asset)
	}

	start = end
	binary.LittleEndian.PutUint32(bs[start:], account.nonce)

	return bs
}

func decodeAccount(bs []byte) *account_T {
	var start int
	end := 8
	account := &account_T{}
	account.balance = binary.LittleEndian.Uint64(bs[:end])

	account.assets = make(map[string]uint64)
	length := (len(bs) - 12) / 40

	for i := 0; i < length; i++ {
		start = end
		end += 40
		account.assets[hex.EncodeToString(bs[start:32])] = binary.LittleEndian.Uint64(bs[start + 32:end])
	}

	start = end

	account.nonce = binary.LittleEndian.Uint32(bs[start:])


	return account
}

func (account *account_T)hash() [32]byte {
	return sha256.Sum256(account.encode())
}
