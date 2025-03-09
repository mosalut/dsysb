// dsysb 

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"net/http"
	"log"
)

const dsysbId = "0000000000000000000000000000000000000000000000000000000000000000"

type state_T struct {
	prevHash [36]byte
	bits uint32
	assets assetPool_T
	accounts map[string]*account_T
}

func (state *state_T)encode() []byte {
	assetLength := len(state.assets) * asset_length

	var accountLength int
	for _, account := range state.accounts {
		accountLength += 14 + len(account.assets) * 40 + 34 // 14 = 8 + 4 + 2, 40 = 32 + 8
	}

	length := 44 + assetLength + accountLength // 44 = 36 + 4 + 4
	bs := make([]byte, length, length)
	var start int
	end := 36
	copy(bs[:end], state.prevHash[:])

	start = end
	end += 4
	binary.LittleEndian.PutUint32(bs[start:end], state.bits)

	start = end
	end += assetLength
	copy(bs[start:end], state.assets.encode())

	for k, account := range state.accounts {
		start = end
		end += 34
		copy(bs[start:end], k)

		accountBytes := account.encode()
		start = end
		end += len(accountBytes)
		copy(bs[start:end], accountBytes)

		start = end
		end += 2
		binary.LittleEndian.PutUint16(bs[start:end], uint16(len(account.assets)))
	}

	start = end

	binary.LittleEndian.PutUint32(bs[start:], uint32(accountLength))
	return bs
}

func decodeState(bs []byte) *state_T {
	var start int
	end := 36

	state := &state_T{}
	copy(state.prevHash[:], bs[:end])

	start = end
	end += 4
	state.bits = binary.LittleEndian.Uint32(bs[start:end])


	start = len(bs) - 4
	accountBytesLength := int(binary.LittleEndian.Uint32(bs[start:]))
	assetEndPosition := len(bs) - accountBytesLength - 4

	start = end
	state.assets = decodeAssetPool(bs[start:assetEndPosition])

	state.accounts = make(map[string]*account_T)

	start = len(bs) - 4

	var assetsInAccount int

	for start > assetEndPosition {
		end = start
		start = end - 2

		assetsInAccount = int(binary.LittleEndian.Uint16(bs[start:end]))

		end = start
		start -= 12 + assetsInAccount * 40

		accountBytes := bs[start:end]

		end = start
		start -= 34
		state.accounts[string(bs[start:end])] = decodeAccount(accountBytes)
	}

	return state
}

func (state *state_T)hash() [32]byte {
	return sha256.Sum256(state.encode())
}

func getState() *state_T {
	bs, err := chainDB.Get([]byte("state"), nil)
	if err != nil {
		print(log_error, err, "state")
		log.Fatal(err, "`state` data has been broken.")
	}

	return decodeState(bs)
}

func (state *state_T)update() error {
	err := chainDB.Put([]byte("state"), state.encode(), nil)
	if err != nil {
		return err
	}

	return nil
}

func initState() {
	_, err := chainDB.Get([]byte("state"), nil)
	if err != nil {
		state := &state_T{}
		state.bits = binary.LittleEndian.Uint32(difficult_1_target[:])
		state.assets = make(assetPool_T)
		state.accounts = make(map[string]*account_T)
		err := state.update()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func stateHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	state := getState()

	stateBytes := state.encode()

	writeResult(w, responseResult_T{true, "ok", stateBytes})
}
