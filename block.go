// dsysb

package main

import (
	"strconv"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"fmt"
)

const (
	bh_length = 152
	bh_hash_position = 36
	bh_state_root_position = 72
	bh_transaction_root_position = 104
	bh_bits_position = 136
	bh_time_position = 140
	bh_nonce_position = 148
)

// Head:
// 	[:36] - prev hash
// 	[36:72] - hash
//	[72:104] - state root
//	[104:136] - transaction root
// 	[136:140] - bits
// 	[140:148] - timestamp
// 	[148:152] - nonce
type blockHead_T struct {
	PrevHash [36]byte `json:"prevHash"`
	Hash [36]byte `json:"hash"`
	StateRoot [32]byte `json:"stateRoot"`
	TransactionRoot [32]byte `json:"transactionRoot"`
	Bits [4]byte `json:"bits"`
	Timestamp [8]byte `json:"timestamp"`
	Nonce [4]byte `json:"nonce"`
}

func (head blockHead_T) String() string {
	return "head:" +
	"\n\tprev hash:" + hex.EncodeToString(head.PrevHash[:]) +
	"\n\thash:" + hex.EncodeToString(head.Hash[:]) +
	"\n\tstate root:" + hex.EncodeToString(head.StateRoot[:]) +
	"\n\ttransaction root:" + hex.EncodeToString(head.TransactionRoot[:]) +
	"\n\tbits:" + hex.EncodeToString(head.Bits[:]) +
	"\n\ttimestamp:" + fmt.Sprintf("%d", binary.LittleEndian.Uint64(head.Timestamp[:])) +
	"\n\tnonce:" + hex.EncodeToString(head.Nonce[:])
}

func (head *blockHead_T) encode () []byte {
	bs := make([]byte, bh_length, bh_length)
	copy(bs[:bh_hash_position], head.PrevHash[:])
	copy(bs[bh_hash_position:bh_state_root_position], head.Hash[:])
	copy(bs[bh_state_root_position:bh_transaction_root_position], head.StateRoot[:])
	copy(bs[bh_transaction_root_position:bh_bits_position], head.TransactionRoot[:])
	copy(bs[bh_bits_position:bh_time_position], head.Bits[:])
	copy(bs[bh_time_position:bh_nonce_position], head.Timestamp[:])
	copy(bs[bh_nonce_position:], head.Nonce[:])

	return bs
}

func decodeBlockHead(bs []byte) *blockHead_T {
	return &blockHead_T {
		[36]byte(bs[:bh_hash_position]),
		[36]byte(bs[bh_hash_position:bh_state_root_position]),
		[32]byte(bs[bh_state_root_position:bh_transaction_root_position]),
		[32]byte(bs[bh_transaction_root_position:bh_bits_position]),
		[4]byte(bs[bh_bits_position:bh_time_position]),
		[8]byte(bs[bh_time_position:bh_nonce_position]),
		[4]byte(bs[bh_nonce_position:]),
	}
}

type blockBody_T struct {
	Transactions txPool_T `json:"transactions"`
}

func (body *blockBody_T) encode () []byte {
	bs := body.Transactions.encode()

	return bs
}

func decodeBlockBody(bs []byte) *blockBody_T {
	body := &blockBody_T{}
	body.Transactions = decodeTxPool(bs)

	return body
}

type block_T struct {
	Head *blockHead_T `json:"head"`
	Body *blockBody_T `json:"body"`
}

func (block *block_T) encode() []byte {
	return append(block.Head.encode(), block.Body.encode()...)
}

func decodeBlock(bs []byte) *block_T {
	block := &block_T{}
	block.Head = decodeBlockHead(bs[:bh_length])
	block.Body = decodeBlockBody(bs[bh_length:])

	return block
}

const genesisPrevHash = "000000000000000000000000000000000000000000000000000000000000000000000000" // [36]byte{}

func getHashBlock() (*block_T, error) {
	state := getState()
	blockBytes, err := chainDB.Get(state.prevHash[32:], nil)
	if err != nil {
		print(log_error, err, fmt.Sprintf("%072x", state.prevHash))
		return nil, err
	}

	return decodeBlock(blockBytes), nil
}

func getBlock(hashBytes []byte) (*block_T, error) {
	blockBytes, err := chainDB.Get(hashBytes, nil)
	if err != nil {
		print(log_error, err)
		return nil, err
	}

	block := &block_T {}
	block.Head = decodeBlockHead(blockBytes[:bh_length])
	block.Body = decodeBlockBody(blockBytes[bh_length:])

	return block, nil
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
		blockchain = append(blockchain, block.Head)

		if hex.EncodeToString(block.Head.PrevHash[:]) == genesisPrevHash {
			break
		}
		block, err = getBlock(block.Head.PrevHash[32:])
		if err != nil {
			writeResult(w, responseResult_T{false, err.Error(), nil})
			return
		}
	}

	writeResult(w, responseResult_T{true, "ok", blockchain.encode()})
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
	hash := values.Get("index")

	height, err := strconv.Atoi(hash)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error() + " height should be a number!", nil})
		return
	}

	buffer := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(buffer, uint32(height))

	block, err := getBlock(buffer)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error() + " height should be a number!", nil})
		return
	}

	writeResult(w, responseResult_T{true, "ok", block.encode()})
}
