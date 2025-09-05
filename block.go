// dsysb

package main

import (
	"strconv"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"time"
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

func decodeBlockBody(bs []byte) (*blockBody_T, error) {
	body := &blockBody_T{}
	var err error
	body.transactions, err = decodeTxPool(bs)
	if err != nil {
		return nil, err
	}

	return body, nil
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

func decodeBlock(bs []byte) (*block_T, error) {
	var err error
	length := len(bs)
	block := &block_T{}
	block.statePosition = binary.LittleEndian.Uint32(bs[length - 4:])
	block.head = decodeBlockHead(bs[:bh_length])
	block.body, err = decodeBlockBody(bs[bh_length:int(block.statePosition)])
	if err != nil {
		return nil, err
	}
	block.state = decodeState(bs[int(block.statePosition):length - 4])

	return block, nil
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

	hash0 := hex.EncodeToString(hash)
	hash1 := hex.EncodeToString(block.head.hash[:])
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
	block.body, err = decodeBlockBody(blockBytes[bh_length:int(statePosition)])
	if err != nil {
		return nil, err
	}
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

	index := binary.LittleEndian.Uint32(block.head.hash[32:])
	coinbase := &coinbase_T {}
	coinbase.rewards(index)

	for _, tx := range block.body.transactions[1:] {
		err = tx.validate(block.head, true)
		if err != nil{
			return err
		}
	}

	block.body.transactions[0].(*coinbase_T).amount = coinbase.amount
	block.body.transactions[0].(*coinbase_T).nonce = index - 1
	*block.state = *prevBlock.state
	block.state.count(block.body.transactions[0].(*coinbase_T))

	for k, tx := range block.body.transactions[1:] {
		err = tx.count(block.state, block.body.transactions[0].(*coinbase_T), k)
		if err != nil{
			return err
		}
	}

	block.body.transactions[0].count(block.state, nil, 0)

	if hex.EncodeToString(newMerkleTree(block.body.transactions).data[:]) != hex.EncodeToString(block.head.transactionRoot[:]) {
		return errTransactionRootHash
	}

	stateHashB := block.state.hash()
	stateHashS := hex.EncodeToString(stateHashB[:])
	if stateHashS != hex.EncodeToString(block.head.stateRoot[:]) {
		return errStateRoot
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

func makeBlockForMine(address string) (*block_T, error) {
	block, err := getHashBlock()
	if err == errZeroBlock {
		print(log_warning, err)
	} else if err != nil {
		return nil, err
	}

	err = adjustTarget(block)
	if err != nil {
		return nil, err
	}

	index := binary.LittleEndian.Uint32(block.head.hash[32:])
	coinbase := &coinbase_T {}
	coinbase.to = address
	coinbase.rewards(index + 1)
	coinbase.nonce = index

	block.body.transactions = txPool_T{ coinbase }
	var txs txPool_T

	if len(transactionPool) <= 511 {
		txs = transactionPool
	//	block.body.transactions = append(block.body.transactions, transactionPool...)
		transactionPool = transactionPool[:0]
	} else {
		txs = transactionPool[:511]
	//	block.body.transactions = append(block.body.transactions, transactionPool[:511]...)
		transactionPool = transactionPool[511:]
	}

	for _, tx := range txs {
		err := tx.validate(block.head, true)
		if err != nil{
			err605 := makeBlockError{"605:", err.Error()}
			print(log_warning, err605)
			noticeErrorBroadcast(err605)

			continue
		}

		block.body.transactions = append(block.body.transactions, tx)
	}

	block.state.count(coinbase)

	for i := 1; i < len(block.body.transactions); {
		stateBack := block.state.encode()
		err := block.body.transactions[i].count(block.state, block.body.transactions[0].(*coinbase_T), i)
		if err != nil {
			if len(block.body.transactions) == i + 1 {
				block.body.transactions = block.body.transactions[:i]
			} else {
				block.body.transactions = append(block.body.transactions[:i], block.body.transactions[i + 1:]...)
			}
			block.state = decodeState(stateBack)

			err605 := makeBlockError{"605:", err.Error()}
			print(log_warning, err605)
			noticeErrorBroadcast(err605)

			continue
		}

		i++
	}

	block.body.transactions[0].count(block.state, nil, 0)

	block.head.prevHash = block.head.hash
	block.head.stateRoot = block.state.hash()
	block.head.transactionRoot = newMerkleTree(block.body.transactions).data
	binary.LittleEndian.PutUint64(block.head.timestamp[:], uint64(time.Now().Unix()))

	return block, nil
}
