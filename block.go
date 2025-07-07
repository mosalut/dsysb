// dsysb

package main

import (
	"strconv"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"errors"
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
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

var errBlockIdNotMatch = errors.New("block hash and index are not match.")
var errBlockHashFormat = errors.New("invalid block hash format.")
var errPrevHashNotMatch = errors.New("state prev hash and block prev hash are not match.")
var errZeroBlock = errors.New("Zero block")
var errBlockHashing = errors.New("block hash and it's data are not match")
var errBits = errors.New("The bits are not match")
var errTransactionRootHash = errors.New("The transaction root and it's hash root are not match")
var errStateRootHash = errors.New("The state and it's hash root are not match")
var errStateRoot = errors.New("The state roots are not match")

// Head:
// 	[:36] - prev hash
// 	[36:72] - hash
//	[72:104] - state root
//	[104:136] - transaction root
// 	[136:140] - bits
// 	[140:148] - timestamp
// 	[148:152] - nonce
type blockHead_T struct {
	prevHash [36]byte
	hash [36]byte
	stateRoot [32]byte
	transactionRoot [32]byte
	bits [4]byte
	timestamp [8]byte
	nonce [4]byte
}

func (head blockHead_T) String() string {
	return "head:" +
	"\n\tprev hash:" + hex.EncodeToString(head.prevHash[:]) +
	"\n\thash:" + hex.EncodeToString(head.hash[:]) +
	"\n\tstate root:" + hex.EncodeToString(head.stateRoot[:]) +
	"\n\ttransaction root:" + hex.EncodeToString(head.transactionRoot[:]) +
	"\n\tbits:" + hex.EncodeToString(head.bits[:]) +
	"\n\ttimestamp:" + fmt.Sprintf("%d", binary.LittleEndian.Uint64(head.timestamp[:])) +
	"\n\tnonce:" + hex.EncodeToString(head.nonce[:])
}

func (head *blockHead_T) encode () []byte {
	bs := make([]byte, bh_length, bh_length)
	copy(bs[:bh_hash_position], head.prevHash[:])
	copy(bs[bh_hash_position:bh_state_root_position], head.hash[:])
	copy(bs[bh_state_root_position:bh_transaction_root_position], head.stateRoot[:])
	copy(bs[bh_transaction_root_position:bh_bits_position], head.transactionRoot[:])
	copy(bs[bh_bits_position:bh_time_position], head.bits[:])
	copy(bs[bh_time_position:bh_nonce_position], head.timestamp[:])
	copy(bs[bh_nonce_position:], head.nonce[:])

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

func (head *blockHead_T) hashing() [32]byte {
	bs := make([]byte, 116, 116)
	copy(bs[:36], head.prevHash[:])
	copy(bs[36:68], head.stateRoot[:])
	copy(bs[68:100], head.transactionRoot[:])
	copy(bs[100:104], head.bits[:])
	copy(bs[104:112], head.timestamp[:])
	copy(bs[112:], head.nonce[:])

	return sha256.Sum256(bs)
}

type blockBody_T struct {
	transactions txPool_T
}

func (body *blockBody_T) encode () []byte {
	bs := body.transactions.encode()

	return bs
}

func decodeBlockBody(bs []byte) *blockBody_T {
	body := &blockBody_T{}
	body.transactions = decodeTxPool(bs)

	return body
}

func (body *blockBody_T) String() string {
	s := "body:\n"
	for _, tx := range body.transactions {
		s += tx.String()
	}

	return s
}

type block_T struct {
	head *blockHead_T
	body *blockBody_T
	state *state_T
	statePosition uint32
}

func (block *block_T) encode() []byte {
	bs := append(block.head.encode(), block.body.encode()...)
	statePosition := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(statePosition, uint32(len(bs)))
	bs = append(bs, block.state.encode()...)
	bs = append(bs, statePosition...)

	return bs
}

func decodeBlock(bs []byte) *block_T {
	length := len(bs)
	block := &block_T{}
	block.statePosition = binary.LittleEndian.Uint32(bs[length - 4:])
	block.head = decodeBlockHead(bs[:bh_length])
	block.body = decodeBlockBody(bs[bh_length:int(block.statePosition)])
	block.state = decodeState(bs[int(block.statePosition):length - 4])

	return block
}

func (block *block_T) String() string {
	return block.head.String() + "\n" + block.body.String() + "\n" + block.state.String()
}

const genesisPrevHash = "000000000000000000000000000000000000000000000000000000000000000000000000" // [36]byte{}

func getHashBlock() (*block_T, error) {
	indexB, err := chainDB.Get([]byte("index"), nil)
	if err != nil {
		return nil, err
	}

	index := binary.LittleEndian.Uint32(indexB)
	if index == 0 {
		return nil, errZeroBlock
	}

	block, err := getBlock(indexB)
	if err != nil {
		print(log_error, err)
		return nil, err
	}

	return block, nil
}

func getBlockByHash(hash []byte) (*block_T, error) {
	if len(hash) != 36 {
		return nil, errBlockHashFormat
	}

	block, err := getBlock(hash[32:])
	if err != nil {
		return nil, err
	}

	hash0 := fmt.Sprintf("%x", hash)
	hash1 := fmt.Sprintf("%x", block.head.hash)
	if hash0 != hash1 {
		return nil, errBlockIdNotMatch
	}

	return block, nil
}

func getBlock(hashBytes []byte) (*block_T, error) {
	blockBytes, err := chainDB.Get(hashBytes, nil)
	if err != nil {
		return nil, err
	}

	length := len(blockBytes)
	statePosition := binary.LittleEndian.Uint32(blockBytes[length - 4:length])

	block := &block_T {}
	block.head = decodeBlockHead(blockBytes[:bh_length])
	block.body = decodeBlockBody(blockBytes[bh_length:int(statePosition)])
	block.state = decodeState(blockBytes[int(statePosition):length - 4])

	return block, nil
}

func (block *block_T)Append() error {
	err := block.validate()
	if err != nil {
		return err
	}

	batch := &leveldb.Batch{}
	batch.Put([]byte("index"), block.head.hash[32:])
	batch.Put(block.head.hash[32:], block.encode())

	err = chainDB.Write(batch, nil)
	if err != nil {
		return err
	}

	noticeAppendBroadcast(block)

	return nil
}

func (block *block_T)validate() error {
	prevBlock, err := getHashBlock()
	if err != nil {
		return err
	}

	prevHash := hex.EncodeToString(prevBlock.head.hash[:])
	blockPrevHash := hex.EncodeToString(block.head.prevHash[:])

	if prevHash != blockPrevHash {
		return errPrevHashNotMatch
	}

	hing := block.head.hashing()
	h1 := hex.EncodeToString(hing[:])
	h2 := hex.EncodeToString(block.head.hash[:32])
	if h1 != h2 {
		return errBlockHashing
	}

	err = adjustTarget(prevBlock)
	if err != nil {
		return err
	}

	bits0 := binary.LittleEndian.Uint32(prevBlock.head.bits[:])
	bits1 := binary.LittleEndian.Uint32(block.head.bits[:])

	if bits0 != bits1 {
		return errBits
	}

	if hex.EncodeToString(newMerkleTree(block.body.transactions).data[:]) != hex.EncodeToString(block.head.transactionRoot[:]) {
		return errTransactionRootHash
	}

	stateHashB := block.state.hash()
	stateHashS := hex.EncodeToString(stateHashB[:])
	if stateHashS != hex.EncodeToString(block.head.stateRoot[:]) {
		return errStateRoot
	}

	for k, tx := range block.body.transactions[1:] {
		err = tx.validate(true)
		if err != nil{
			return err
		}

		err = tx.count(prevBlock.state, block.body.transactions[0].(*coinbase_T), k)
		if err != nil{
			return err
		}
	}

	err = block.body.transactions[0].validate(true)
	if err != nil{
		return err
	}

	err = block.body.transactions[0].count(prevBlock.state, nil, 0)
	if err != nil{
		return err
	}

	return nil
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
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	writeResult(w, responseResult_T{true, "ok", block.encode()})
}

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
