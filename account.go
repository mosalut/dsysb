// dsysb

package main

import (
	"sort"
	"math/big"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
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

	astLength := len(account.assets)
	keys := make([]string, 0, astLength)
	for k, _ := range account.assets {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		a := big.NewInt(0)
		a.SetString(keys[i], 16)

		b := big.NewInt(0)
		b.SetString(keys[j], 16)

		return a.Cmp(b) > 0
	})

	for _, k := range keys {
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
		binary.LittleEndian.PutUint64(bs[start:end], account.assets[k])
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
		account.assets[hex.EncodeToString(bs[start:start + 32])] = binary.LittleEndian.Uint64(bs[start + 32:end])
	}

	start = end

	account.nonce = binary.LittleEndian.Uint32(bs[start:])


	return account
}

func (account *account_T)hash() [32]byte {
	return sha256.Sum256(account.encode())
}

func (account *account_T)String() string {
	value := fmt.Sprintf("\tbalance: %d", account.balance)
	value += "\tassets:\n"
	for k, asset := range account.assets {
		value += fmt.Sprintf("\t\t%s: %v\n", k, asset)
	}
	value += fmt.Sprintf("\tnonce: %d\n", account.nonce)

	return value
}
