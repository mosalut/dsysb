// dsysb

package main

import (
	"encoding/hex"
	"net/http"
)

func getTransaction2Handler(w http.ResponseWriter, req *http.Request) {
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
	txid := values.Get("txid")

	block, err := getHashBlock()
	if err != nil {
		writeResult2(w, responseResult2_T{false, err.Error(), nil})
		return
	}

	for _, tx := range block.body.transactions {
		h := tx.hash()
		if txid == hex.EncodeToString(h[:]) {
			writeResult2(w, responseResult2_T{true, "ok", tx.Map()})
			return
		}
	}

	for block != nil {
		if hex.EncodeToString(block.head.prevHash[:]) == genesisPrevHash {
			break
		}

		block, err = getBlock(block.head.prevHash[32:])
		if err != nil {
			writeResult2(w, responseResult2_T{false, err.Error(), nil})
			return
		}

		for _, tx := range block.body.transactions {
			h := tx.hash()
			if txid == hex.EncodeToString(h[:]) {
				writeResult2(w, responseResult2_T{true, "ok", tx.Map()})
				return
			}
		}
	}

	writeResult2(w, responseResult2_T{false, "Not found the txid:" + txid, nil})
}
