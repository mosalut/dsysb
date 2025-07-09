// dsysb

package main

import (
	"strconv"
	"encoding/binary"
	"encoding/hex"
	"net/http"
)

func block2Handler(w http.ResponseWriter, req *http.Request) {
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
	hash := values.Get("index")

	height, err := strconv.Atoi(hash)
	if err != nil {
		writeResult2(w, responseResult2_T{false, err.Error() + " height should be a number!", nil})
		return
	}

	buffer := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(height))

	block, err := getBlock(buffer)
	if err != nil {
		writeResult2(w, responseResult2_T{false, err.Error(), nil})
		return
	}

	block2 := struct {
		Head *struct {
			PrevHash string `json:"prevHash"`
			Hash string `json:"hash"`
			StateRoot string `json:"stateRoot"`
			TransactionRoot string `json:"transactionRoot"`
			Bits string `json:"bits"`
			Timestamp int64 `json:"timestamp"`
			Nonce uint32 `json:"nonce"`
		} `json:"head"`
		Body []string `json:"body"`
	} {}

	block2.Head = &struct {
		PrevHash string `json:"prevHash"`
		Hash string `json:"hash"`
		StateRoot string `json:"stateRoot"`
		TransactionRoot string `json:"transactionRoot"`
		Bits string `json:"bits"`
		Timestamp int64 `json:"timestamp"`
		Nonce uint32 `json:"nonce"`
	} {}
	block2.Head.PrevHash = hex.EncodeToString(block.head.prevHash[:])
	block2.Head.Hash = hex.EncodeToString(block.head.hash[:])
	block2.Head.StateRoot = hex.EncodeToString(block.head.stateRoot[:])
	block2.Head.TransactionRoot = hex.EncodeToString(block.head.transactionRoot[:])
	block2.Head.Bits = hex.EncodeToString(block.head.bits[:])
	block2.Head.Timestamp = int64(binary.LittleEndian.Uint64(block.head.timestamp[:]))
	block2.Head.Nonce = uint32(binary.LittleEndian.Uint32(block.head.nonce[:]))

	tLength := len(block.body.transactions)
	block2.Body = make([]string, tLength, tLength)

	for k, tx := range block.body.transactions {
		h := tx.hash()
		block2.Body[k] = hex.EncodeToString(h[:])
	}

	writeResult2(w, responseResult2_T{true, "ok", block2})
}
