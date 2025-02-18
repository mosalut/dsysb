// dsysb

package main

import (
	"strconv"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"fmt"
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
	bs := make([]byte, 152, 152)
	copy(bs[:36], head.PrevHash[:])
	copy(bs[36:72], head.Hash[:])
	copy(bs[72:104], head.StateRoot[:])
	copy(bs[104:136], head.TransactionRoot[:])
	copy(bs[136:140], head.Bits[:])
	copy(bs[140:148], head.Timestamp[:])
	copy(bs[148:], head.Nonce[:])

	return bs
}

func decodeBlockHead(bs []byte) *blockHead_T {
	return &blockHead_T {
		[36]byte(bs[:36]),
		[36]byte(bs[36:72]),
		[32]byte(bs[72:104]),
		[32]byte(bs[104:136]),
		[4]byte(bs[136:140]),
		[8]byte(bs[140:148]),
		[4]byte(bs[148:]),
	}
}

type blockBody_T struct {
	Transactions []*transaction_T `json:"transactions"`
}

func (body *blockBody_T) encode () []byte {
	bs, err := json.Marshal(body)
	if err != nil {
		print(log_error, err)
		return nil
	}

	return bs
}

func decodeBlockBody(bs []byte) *blockBody_T {
	body := &blockBody_T{}
	err := json.Unmarshal(bs, body)
	if err != nil {
		print(log_error, err)
		return nil
	}

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
	block.Head = decodeBlockHead(bs[:152])
	block.Body = decodeBlockBody(bs[152:])

	return block
}

const genesisPrevHash = "000000000000000000000000000000000000000000000000000000000000000000000000" // [36]byte{}

func getHashBlock() *block_T {
	state := getState()
	blockBytes, err := chainDB.Get(state.PrevHash[32:], nil)
	if err != nil {
		print(log_error, err, fmt.Sprintf("%072x", state.PrevHash))
		return nil
	}

	return decodeBlock(blockBytes)
}

func getBlock(hashBytes []byte) *block_T {
	blockBytes, err := chainDB.Get(hashBytes, nil)
	if err != nil {
		print(log_error, err)
		return nil
	}

	block := &block_T {}
	block.Head = decodeBlockHead(blockBytes[:152])
	block.Body = decodeBlockBody(blockBytes[152:])

	return block
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

	blockchain := make([]*block_T, 0, number)
	block := getHashBlock()
	for i := 0; i < number && block != nil; i++ {
		blockchain = append(blockchain, block)

		if hex.EncodeToString(block.Head.PrevHash[:]) == genesisPrevHash {
			break
		}
		block = getBlock(block.Head.PrevHash[32:])
	}

	jsonData, err := json.Marshal(blockchain)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	writeResult(w, responseResult_T{true, "ok", jsonData})
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

	block := getBlock(buffer)

	writeResult(w, responseResult_T{true, "ok", block.encode()})
}
