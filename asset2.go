// dsysb

package main

import (
	"encoding/hex"
	"net/http"
)

type asset_J struct {
	Name string `json:"name"`
	Symbol string `json:"symbol"`
	Decimals uint8 `json:"decimals"`
	TotalSupply uint64 `json:"totalSupply"`
	Price uint32 `json:"price"`
	blocks uint32 `json:"blocks"`
	Remain uint32 `json:"remain"`
}

type assetPool_J map[string]*asset_J  // key is an assetId

func assets2Handler(w http.ResponseWriter, req *http.Request) {
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

	writeResult2(w, responseResult2_T{true, "ok", assets})
}

func asset2Handler(w http.ResponseWriter, req *http.Request) {
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
	assetId := values.Get("id")

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult2(w, responseResult2_T{false, "dsysb inner error", nil})
		return
	}

	for _, asset := range state.assets {
		h := asset.hash()
		aId := hex.EncodeToString(h[:])

		if aId == assetId {
			assetJ := &asset_J {
				asset.name,
				asset.symbol,
				asset.decimals,
				asset.totalSupply,
				asset.price,
				asset.blocks,
				asset.remain,
			}
			writeResult2(w, responseResult2_T{true, "ok", assetJ})
			return
		}
	}

	writeResult2(w, responseResult2_T{false, "asset " + assetId + " does not exist", nil})
}
