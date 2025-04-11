// dsysb

package main

import (
	"math"
	"strconv"
	"net"
	"net/http"
	"encoding/hex"
	"encoding/binary"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

// 1f00ffff
// const firstBlock = "00000000000000000000000000000000000000000000000000000000000000000000000000004062db4e48095bb6c67e5f596e5a6dd7bdf1cdb08f768a4484758078eb1f010000005264f7825a364111c70007b13eca64bbe0f0ec30612df8420c1a2e3beefa11e05089f17493f12d7e9ce5f8bfcd2d918c7eafd3b048ea221ee97a62ecb60941d01f00ffff99cef06700000000505f0000004443557a325a32443843395943375a5761564277417831367679677041473466535400743ba40b000000dee5ba3c1f00ffff4443557a325a32443843395943375a5761564277417831367679677041473466535400743ba40b00000000000000000030000000c7000000"

// 1d00ffff
const firstBlock = "00000000000000000000000000000000000000000000000000000000000000000000000000000000f2d0aa718431a47c7c7d4257eb2429f0a61788ca8bf28c06c02e2ddf01000000a51e6a3b08dd85419cf3448f837cd95f98762aefd5f4a789140264758e6a94eec7193cb5d332bed0529d67c8408e650917cb40b546d2ae13cf1eff5fd5b1a11e1d00ffff3e32f7670000000047ded50000443839326b786261765a55556d6a3544486f56434a4146555773574d654a4347515900743ba40b0000003b866fa61d00ffff443839326b786261765a55556d6a3544486f56434a4146555773574d654a4347515900743ba40b00000000000000000030000000c7000000"

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

func sendLastestBlock(rAddr *net.UDPAddr) error {
	block, err := getHashBlock()
	if err != nil {
		return err
	}

	err = transport(rAddr, p2p_add_block_event, block.encode())
	if err != nil {
		return err
	}

	return nil
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

	if number <= 0 || number > math.MaxUint32 {
		writeResult(w, responseResult_T{false, "number must be between 0 and 4294967296", nil})
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
