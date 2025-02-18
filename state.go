// dsysb 

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"net/http"
	"log"
)

type state_T struct {
	PrevHash [36]byte `json:"prevHash"`
	Bits uint32 `json:"bits"`
	Assets map[string]*asset_T `json:"assets"`
	Accounts map[string]*account_T `json:"accounts"`
}

func (state *state_T)encode() []byte {
	bs, err := json.Marshal(state)
	if err != nil {
		print(log_error, err)
		return nil
	}
	return bs
}

func decodeState(bs []byte) *state_T {
	state := &state_T{}
	err := json.Unmarshal(bs, state)
	if err != nil {
		print(log_error, err)
		return nil
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
		state.Bits = binary.LittleEndian.Uint32(difficult_1_target[:])
		state.Assets = make(map[string]*asset_T)
		state.Accounts = make(map[string]*account_T)
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

	stateBytes, err := json.Marshal(state)
	if err != nil {
		print(log_error, err)
		return
	}

	writeResult(w, responseResult_T{true, "ok", stateBytes})
}
