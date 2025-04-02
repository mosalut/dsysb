// dsysb

package main

import (
	"strconv"
	"net"
	"net/http"
	"encoding/hex"
	"encoding/binary"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

const firstBlock = "00000000000000000000000000000000000000000000000000000000000000000000000000000ca4383d9db17acc0ccfd39f876c9d259ed3e71ef582e65da058018e7cca010000002c57e1558dd4b577a2348aae091bf1f8134a591c7db78f50a53b6697fec59d2321982304d8630de8eb5c2a4f22531e04670d34d405b1c4e13723685f8fae83901f00fffff91aed670000000028140000004443557a325a32443843395943375a5761564277417831367679677041473466535400743ba40b0000008f4b42fa1f00ffff4443557a325a32443843395943375a5761564277417831367679677041473466535400743ba40b000000000000000000300000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c7000000"

type blockchainSync_T struct {
	rAddr *net.UDPAddr
	targetIndex uint32
	synchronizing bool
}

func (chainSync *blockchainSync_T) doing (rAddr *net.UDPAddr) {
	blockchainSync.rAddr = rAddr
	blockchainSync.synchronizing = true
}

func (chainSync *blockchainSync_T) over () {
	blockchainSync.rAddr = nil
	blockchainSync.targetIndex = 0
	blockchainSync.synchronizing = false
}

var blockchainSync blockchainSync_T

type blockchain_T []*blockHead_T

func (chain blockchain_T) encode() []byte {
	length := len(chain)
	lenH := length * bh_length
	bs := make([]byte, lenH, lenH)

	for k, head := range chain {
		copy(bs[k * bh_length:(k + 1) * bh_length], head.encode())
	}

	return bs
}

func decodeBlockchain(bs []byte) blockchain_T {
	length := len(bs)
	if length % bh_length != 0 {
		log.Fatal("Wrong length of block heads")
	}

	lenH := length / bh_length
	blockchain := make([]*blockHead_T, lenH, lenH)
	for i := 0; i < lenH; i++ {
		head := decodeBlockHead(bs[i * bh_length:(i + 1) * bh_length])
		blockchain[i] = head
	}

	return blockchain
}

func rollbackChain(startIndex uint32) error {
	height, err := getIndex()
	if err != nil {
		return err
	}

	batch := &leveldb.Batch{}
	indexB := make([]byte, 4, 4)
	var index uint32
	for index = height; index > startIndex; index-- {
		binary.LittleEndian.PutUint32(indexB, index)
		batch.Delete(indexB)
	}
	binary.LittleEndian.PutUint32(indexB, index)
	batch.Put([]byte("index"), indexB)
	err = chainDB.Write(batch, nil)
	if err != nil {
		blockchainSync.over()
		return err
	}

	return nil
}

func initIndex() {
	/* keepit */
	_, err := chainDB.Get([]byte("index"), nil)
	if err == leveldb.ErrNotFound {
		bs, err := hex.DecodeString(firstBlock)
		if err != nil {
			log.Fatal(err)
		}

		batch := &leveldb.Batch{}
		batch.Put([]byte("index"), []byte{1, 0, 0, 0})
		batch.Put([]byte{1, 0, 0, 0}, bs)

		err = chainDB.Write(batch, nil)
		if err != nil {
			log.Fatal(err)
		}
	}
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

	blockchain := make(blockchain_T, 0, number)
	block, err := getHashBlock()
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	for i := 0; i < number && block != nil; i++ {
		blockchain = append(blockchain, block.head)

		if hex.EncodeToString(block.head.prevHash[:]) == genesisPrevHash {
			break
		}
		block, err = getBlock(block.head.prevHash[32:])
		if err != nil {
			writeResult(w, responseResult_T{false, err.Error(), nil})
			return
		}
	}

	writeResult(w, responseResult_T{true, "ok", blockchain.encode()})
}
