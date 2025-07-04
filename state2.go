// dsysb 

package main

import (
	"encoding/hex"
	"net/http"
)

type state_J struct {
	Assets assetPool_J `json:"assets"`
	Accounts map[string]*account_J `json:"accounts"`
	Tasks taskPool_J `json:"tasks"`
}

func state2Handler(w http.ResponseWriter, req *http.Request) {
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

	assets := make(assetPool_J)
	for k, asset := range state.assets {
		assets[k] = &asset_J {
			asset.name,
			asset.symbol,
			asset.decimals,
			asset.totalSupply,
			asset.price,
			asset.blocks,
			asset.remain,
		}
	}

	accounts := make(map[string]*account_J)
	for k, account := range state.accounts {
		accounts[k] = &account_J {
			account.balance,
			account.assets,
			account.nonce,
		}
	}

	tLength := len(state.tasks)
	tasks := make(taskPool_J, tLength, tLength)
	for k, task := range state.tasks {
		tasks[k] = &task_J {
			task.address,
			hex.EncodeToString(task.instructs[:]),
			hex.EncodeToString(task.vData[:]),
		}
	}

	stateJ := &state_J {
		assets,
		accounts,
		tasks,
	}

	writeResult2(w, responseResult2_T{true, "ok", stateJ})
}
