package main

import (
	"net/http"
)

type blockHead_T struct {
	Index int `json:"index"`
	Timestamp int64 `json:"timestamp"`
	Hash string `json:"hash"`
	PrevHash string `json:"prev_hash"`
	Nonce int `json:"nonce"`
}

type blockBody_T struct {
	Transactions []*transaction_T `json:"transactions"`
}

type block_T struct {
	Head *blockHead_T `json:"head"`
	Body *blockBody_T `json:"body"`
}

var blockchain = make([]*block_T, 0, 2048)

func getTillBlock() *block_T {
	if len(blockchain) == 0 {
		return nil
	}

	return blockchain[len(blockchain) - 1]
}

func tillBlockHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	block := getTillBlock()
	writeResult(w, responseResult_T{true, "ok", block})
}

func blockchainHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	writeResult(w, responseResult_T{true, "ok", blockchain})
}
