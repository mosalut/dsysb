// dsysb

package main

import (
	"math"
	"strconv"
	"net"
	"net/http"
	"encoding/hex"
	"encoding/binary"
	"time"
	"log"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

// 1d00ffff mainnet
const firstBlock = "00000000000000000000000000000000000000000000000000000000000000000000000000000000ba92bac6b89f318294ec697020b6fecc2ef53883b554649a7e13872101000000d35dcb68f6cbb9f34f4a0fa5f82ace01314c46817a0b5b143918873f33c23c5dec8259ef6bd2e150449e209cc37987c9641a937f9f466580b71be4c7909289861d00ffff9a5c13690000000008dfa3cf000000002f00004435565659364157636f584c394251597a577458567858634c5671446b626343646d00743ba40b000000000000004435565659364157636f584c394251597a577458567858634c5671446b626343646d00743ba40b0000000000000000003000000000000000cd000000"


// 1f00ffff devnet
// const firstBlock = "0000000000000000000000000000000000000000000000000000000000000000000000000000a242f8c25d1baf144068aefefdb249c95f04ea34cfae1117d2a9d261ac3a010000008667b8391835267215fd9e86bf8e6dd036f306b228275982065022c59522b09f47cedcc960bd5050a2a2645a04b3da3926d0b3a42e82403385cff37bce8b39341f00ffffde1d80690000000026400000000000002f00004443557a325a32443843395943375a5761564277417831367679677041473466535400743ba40b000000000000004443557a325a32443843395943375a5761564277417831367679677041473466535400743ba40b0000000000000000003000000000000000cd000000"

// 1f00ffff testnet 2
// const firstBlock = ""

type blockchainSync_T struct {
	rAddr *net.UDPAddr
	synchronizing bool
	targetIndex uint32
	blockIndex uint32
	timestamp int64
}

func (chainSync *blockchainSync_T) doing(rAddr *net.UDPAddr) {
	chainSync.rAddr = rAddr
	chainSync.synchronizing = true
	chainSync.timestamp = time.Now().Unix()
	go func() {
		blockIndex := chainSync.blockIndex
		for {
			if !chainSync.synchronizing {
				return
			}

			time.Sleep(time.Second * 8)
			if blockIndex == chainSync.blockIndex {
				fmt.Println("Request block height again")
				rAddr := chainSync.rAddr
				chainSync.over()
				err := transport(rAddr, p2p_request_height_event, nil)
				if err != nil {
					print(log_info, err)
				}

				return
			}


			blockIndex = chainSync.blockIndex
		}
	}()
}

/*
func (chainSync *blockchainSync_T) synchronizing() bool {
	if chainSync.rAddr == nil {
		return false
	}

	t := time.Now().Unix()
	if t - chainSync.timestamp > 8 {
		chainSync.over()
		return false
	}

	return true
}
*/

func (chainSync *blockchainSync_T) over () {
	chainSync.rAddr = nil
	chainSync.synchronizing = false
	chainSync.targetIndex = 0
	chainSync.timestamp = 0
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

func sendLatestBlock(rAddr *net.UDPAddr) error {
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
			Nonce uint64 `json:"nonce"`
		} `json:"head"`
		Transactions []string `json:"transactions"`
	}, 0, number)

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
				Nonce uint64 `json:"nonce"`
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
			Nonce uint64 `json:"nonce"`
		} {}

		block2.Head.PrevHash = hex.EncodeToString(block.head.prevHash[:])
		block2.Head.Hash = hex.EncodeToString(block.head.hash[:])
		block2.Head.StateRoot = hex.EncodeToString(block.head.stateRoot[:])
		block2.Head.TransactionRoot = hex.EncodeToString(block.head.transactionRoot[:])
		block2.Head.Bits = hex.EncodeToString(block.head.bits[:])
		block2.Head.Timestamp = int64(binary.LittleEndian.Uint64(block.head.timestamp[:]))
		block2.Head.Nonce = block.head.nonce

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
