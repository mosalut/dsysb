// dsysb

package main

import (
	"net/http"
)

type account_J struct {
	Balance uint64 `json:"balance"`
	Assets map[string]uint64 `json:"assets"` // key is an asset id
	Nonce uint32 `json:"nonce"`
}

func account2Handler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	values := req.URL.Query()
	address := values.Get("address")

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult2(w, responseResult2_T{false, "dsysb inner error", nil})
		return
	}

	account, ok := state.accounts[address]
	if !ok {
		writeResult2(w, responseResult2_T{false, "No this account", nil})
	}

	writeResult2(w, responseResult2_T{true, "ok", account})
}

func accounts2Handler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult2(w, responseResult2_T{false, "dsysb inner error", nil})
		return
	}

	accounts := make(map[string]*account_J)
	for k, account := range state.accounts {
		accounts[k] = &account_J {
			account.balance,
			account.assets,
			account.nonce,
		}
	}

	writeResult2(w, responseResult2_T{true, "ok", accounts})
}
