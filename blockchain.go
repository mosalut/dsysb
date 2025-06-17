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
const firstBlock = "0000000000000000000000000000000000000000000000000000000000000000000000000000bfae5225bf5b1b14e3f0860d2a3614e72812a6e98bf894cc4e18e1dd3f57010000008667b8391835267215fd9e86bf8e6dd036f306b228275982065022c59522b09fecb2700c3968bd0481a86f04be9326daa9fb758d7b0db7685c06429ad132b12d1f00ffffef0f27680000000008fe00002e004443557a325a32443843395943375a5761564277417831367679677041473466535400743ba40b000000c56386d54443557a325a32443843395943375a5761564277417831367679677041473466535400743ba40b0000000000000000003000000000000000c8000000"

// 1d00ffff
// const firstBlock = "000000000000000000000000000000000000000000000000000000000000000000000000000000009ac94bf25ec369b1cd22db2d54e3a2e74b7458b769e4120bef46015201000000a6271e23063a9b7a8117afcdcd97ba8ac4cd8506952bc516581aa93309db20c56c2f26d3946c81b68c4000b68ef08b5c2fa3b518c792fd23a44882f51d7138611d00ffff62b1fe6700000000b215bbc6004455417633733741544631733367686b7a504744435453466467535445664a77366700743ba40b000000419d7ac74455417633733741544631733367686b7a504744435453466467535445664a77366700743ba40b00000000000000000030000000c7000000"

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

/* keepfunc */
func initIndex() {
	// keepit
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

func blockchain2Handler(w http.ResponseWriter, req *http.Request) {
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

	blockchain := make([]*struct {
		Head *struct {
			PrevHash string `json:"prevHash"`
			Hash string `json:"hash"`
			StateRoot string `json:"stateRoot"`
			TransactionRoot string `json:"transactionRoot"`
			Bits string `json:"bits"`
			Timestamp int64 `json:"timestamp"`
			Nonce uint32 `json:"nonce"`
		} `json:"head"`
		Transactions []string `json:"transactions"`
	}, 0, number)
	/*
	blockchain := make([]*struct{
	}, 0, number)
	*/

	block, err := getHashBlock()
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	for i := 0; i < number && block != nil; i++ {
		block2 := &struct {
			Head *struct {
				PrevHash string `json:"prevHash"`
				Hash string `json:"hash"`
				StateRoot string `json:"stateRoot"`
				TransactionRoot string `json:"transactionRoot"`
				Bits string `json:"bits"`
				Timestamp int64 `json:"timestamp"`
				Nonce uint32 `json:"nonce"`
			} `json:"head"`
			Transactions []string `json:"transactions"`
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
		block2.Transactions = make([]string, tLength, tLength)
		for k, tx := range block.body.transactions {
			h := tx.hash()
			block2.Transactions[k] = hex.EncodeToString(h[:])
		}

		blockchain = append(blockchain, block2)

		if hex.EncodeToString(block.head.prevHash[:]) == genesisPrevHash {
			break
		}
		block, err = getBlock(block.head.prevHash[32:])
		if err != nil {
			writeResult(w, responseResult_T{false, err.Error(), nil})
			return
		}
	}

	writeResult2(w, responseResult2_T{true, "ok", blockchain})
}
