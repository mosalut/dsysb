package main

import (
	"strconv"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"fmt"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type blockHead_T struct {
	Index uint32 `json:"index"`
	Timestamp int64 `json:"timestamp"`
	Hash string `json:"hash"`
	PrevHash string `json:"prev_hash"`
	Nonce uint32 `json:"nonce"`
}

type blockBody_T struct {
	Transactions []*transaction_T `json:"transactions"`
}

type block_T struct {
	Head *blockHead_T `json:"head"`
	Body *blockBody_T `json:"body"`
}

var chainDB *leveldb.DB
var toolDB *leveldb.DB

func initDB() {
	var err error
	toolDB, err = leveldb.OpenFile("tool.db", nil)
	if err != nil {
		log.Fatal(err)
	}

	chainDB, err = leveldb.OpenFile("chain.db", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func getTillBlock() *block_T {
	tillBytes, err := toolDB.Get([]byte("till"), nil)
	if err != nil {
		print(log_error, err)
		return nil
	}

	fmt.Printf("till: %x\n", tillBytes)

	blockBytes, err := chainDB.Get(tillBytes, nil)
	if err != nil {
		print(log_error, err)
		return nil
	}

	block := block_T{}
	err = json.Unmarshal(blockBytes, &block)
	if err != nil {
		print(log_error, err)
		return nil
	}

	return &block
}

func getBlock(indexBytes []byte) *block_T {
	blockBytes, err := chainDB.Get(indexBytes, nil)
	if err != nil {
		print(log_error, err)
		return nil
	}


	block := block_T{}
	err = json.Unmarshal(blockBytes, &block)
	if err != nil {
		print(log_error, err)
		return nil
	}

	return &block
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

	values := req.URL.Query()
	n := values.Get("n")

	number, err := strconv.Atoi(n)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	var blockchain = make([]*block_T, 0, number)
	block := getTillBlock()
	for i := 0; i < number && block != nil; i++ {
		blockchain = append(blockchain, block)
		if block.Head.PrevHash == "" {
			break
		}
		prevHashBytes, _ := hex.DecodeString(block.Head.PrevHash)
		block = getBlock(prevHashBytes[32:])
	}
	writeResult(w, responseResult_T{true, "ok", blockchain})
}

func blockHandler(w http.ResponseWriter, req *http.Request) {
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
	index := values.Get("index")

	height, err := strconv.Atoi(index)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error() + " height should be a number!", nil})
		return
	}

	buffer := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(height))

	block := getBlock(buffer)

	writeResult(w, responseResult_T{true, "ok", block})
}
