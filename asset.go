// dsysb

package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"encoding/binary"
	"encoding/hex"
	"regexp"
	"net/http"
	"errors"
	"fmt"
)

const (
	asset_length = 40
	asset_symbol_position = 10
	asset_decimals_position = 15
	asset_total_supply_position = 16
	asset_price_position = 24
	asset_blocks_position = 32
	asset_remain_position = 36

	create_asset_length = 210
	create_asset_from_position = 36
	create_asset_nonce_position = 70
	create_asset_fee_position = 74
	create_asset_signer_position = 82
)

const dsysbId = "0000000000000000000000000000000000000000000000000000000000000000"

// The `name` is the asset Name.
// The `symbol` is the asset Symbol.
// `totalSupply` is just total supply.
// `price` represent the Price of the asset that is holden by an block.
// `blocks` represent how many Blocks the asset holden.
type asset_T struct {
	name string
	symbol string
	decimals uint8
	totalSupply uint64
	price uint64
	blocks uint32
	remain uint32
}

func (asset *asset_T) hash() [32]byte {
	bs := make([]byte, asset_length, asset_length)
	copy(bs[:asset_symbol_position], []byte(asset.name))
	copy(bs[asset_symbol_position:asset_decimals_position], []byte(asset.symbol))
	bs[asset_decimals_position] = byte(asset.decimals)
	binary.LittleEndian.PutUint64(bs[asset_total_supply_position:asset_price_position], asset.totalSupply)
	binary.LittleEndian.PutUint64(bs[asset_price_position:asset_blocks_position], asset.price)
	binary.LittleEndian.PutUint32(bs[asset_blocks_position:asset_remain_position], asset.blocks)
	h := sha256.Sum256(bs)

	return h
}

func (asset *asset_T) encode () []byte {
	bs := make([]byte, asset_length, asset_length)
	copy(bs[:asset_symbol_position], []byte(asset.name))
	copy(bs[asset_symbol_position:asset_decimals_position], []byte(asset.symbol))
	bs[asset_decimals_position] = byte(asset.decimals)
	binary.LittleEndian.PutUint64(bs[asset_total_supply_position:asset_price_position], asset.totalSupply)
	binary.LittleEndian.PutUint64(bs[asset_price_position:asset_blocks_position], asset.price)
	binary.LittleEndian.PutUint32(bs[asset_blocks_position:asset_remain_position], asset.blocks)
	binary.LittleEndian.PutUint32(bs[asset_remain_position:], asset.remain)

	return bs
}

func decodeAsset(bs []byte) *asset_T {
	asset := &asset_T{}
	asset.name = string(bs[:asset_symbol_position])
	asset.symbol = string(bs[asset_symbol_position:asset_decimals_position])
	asset.decimals = uint8(bs[asset_decimals_position])
	asset.totalSupply = binary.LittleEndian.Uint64(bs[asset_total_supply_position:asset_price_position])
	asset.price = binary.LittleEndian.Uint64(bs[asset_price_position:asset_blocks_position])
	asset.blocks = binary.LittleEndian.Uint32(bs[asset_blocks_position:asset_remain_position])
	asset.remain = binary.LittleEndian.Uint32(bs[asset_remain_position:])

	return asset
}

func (asset *asset_T) String() string {
	return "\tname:\t" + asset.name +
		"\n\tsymbol:\t" + asset.symbol +
		"\n\tdecimals:\t" + fmt.Sprintf("%d", asset.decimals) +
		"\n\ttotal supply:\t" + fmt.Sprintf("%d Satoshi", asset.totalSupply) +
		"\n\tprice:\t" + fmt.Sprintf("%d Satoshi", asset.price) +
		"\n\tblocks:\t" + fmt.Sprintf("%d", asset.blocks)
}

type assetPool_T map[string]*asset_T  // key is an assetId
func (pool assetPool_T) encode() []byte {
	length := len(pool) * asset_length
	bs := make([]byte, 0, length)
	for _, asset := range pool {
		bs = append(bs, asset.encode()...)
	}

	return bs
}

func decodeAssetPool(bs []byte) assetPool_T {
	pool := make(assetPool_T)
	lengthByte := len(bs)
	if lengthByte == 0 {
		return pool
	}
	length := lengthByte / asset_length
	for i := 0 ; i < length; i++ {
		asset := decodeAsset(bs[i * asset_length:(i + 1) * asset_length])
		h := asset.hash()
		pool[hex.EncodeToString(h[:])] = asset
	}

	return pool
}

// `Blocks` represent how many Blocks the asset holden.
type prolong_T struct {
	assetId [32]byte
	blocks [4]byte
	from [34]byte
}

func assetsHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult(w, responseResult_T{false, "dsysb inner error", nil})
		return
	}

	writeResult(w, responseResult_T{true, "ok", state.assets.encode()})
}

func assetHandler(w http.ResponseWriter, req *http.Request) {
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
	assetId := values.Get("id")

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult(w, responseResult_T{false, "dsysb inner error", nil})
		return
	}

	for _, asset := range state.assets {
		h := asset.hash()
		aId := hex.EncodeToString(h[:])

		if aId == assetId {
			writeResult(w, responseResult_T{true, "ok", asset.encode()})
			return
		}
	}

	writeResult(w, responseResult_T{false, "asset " + assetId + " does not exist", nil})
}

type createAsset_T struct {
	name string
	symbol string
	decimals uint8
	totalSupply uint64
	price uint64
	blocks uint32
	from string
	nonce uint32
	fee uint64
	signer *signer_T
}

func (tx *createAsset_T) hash() [32]byte {
	bs := make([]byte, create_asset_length, create_asset_length)
	copy(bs[:asset_symbol_position], []byte(tx.name))
	copy(bs[asset_symbol_position:asset_decimals_position], []byte(tx.symbol))
	bs[asset_decimals_position] = byte(tx.decimals)
	binary.LittleEndian.PutUint64(bs[asset_total_supply_position:asset_price_position], tx.totalSupply)
	binary.LittleEndian.PutUint64(bs[asset_price_position:asset_blocks_position], tx.price)
	binary.LittleEndian.PutUint32(bs[asset_blocks_position:create_asset_from_position], tx.blocks)
	copy(bs[create_asset_from_position:create_asset_nonce_position], []byte(tx.from))
	binary.LittleEndian.PutUint32(bs[create_asset_nonce_position:create_asset_fee_position], tx.nonce)
	binary.LittleEndian.PutUint64(bs[create_asset_fee_position:create_asset_signer_position], tx.fee)

	return sha256.Sum256(bs)
}

func (ca *createAsset_T) encode() []byte {
	bs := make([]byte, create_asset_length, create_asset_length)
	copy(bs[:asset_symbol_position], []byte(ca.name))
	copy(bs[asset_symbol_position:asset_decimals_position], []byte(ca.symbol))
	bs[asset_decimals_position] = byte(ca.decimals)
	binary.LittleEndian.PutUint64(bs[asset_total_supply_position:asset_price_position], ca.totalSupply)
	binary.LittleEndian.PutUint64(bs[asset_price_position:asset_blocks_position], ca.price)
	binary.LittleEndian.PutUint32(bs[asset_blocks_position:create_asset_from_position], ca.blocks)
	copy(bs[create_asset_from_position:create_asset_nonce_position], []byte(ca.from))
	binary.LittleEndian.PutUint32(bs[create_asset_nonce_position:create_asset_fee_position], ca.nonce)
	binary.LittleEndian.PutUint64(bs[create_asset_fee_position:create_asset_signer_position], ca.fee)
	copy(bs[create_asset_signer_position:], ca.signer.encode())

	return bs
}

func (tx *createAsset_T) encodeForPool() []byte {
	length := create_asset_length + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], create_asset_length)
	copy(bs[2:], tx.encode())

	return bs
}

func decodeCreateAsset(bs []byte) *createAsset_T {
	ca := &createAsset_T{}

	ca.name = string(bytes.Trim(bs[:asset_symbol_position], "\x00 \t\n\r"))
	ca.symbol = string(bytes.Trim(bs[asset_symbol_position:asset_decimals_position], "\x00 \t\n\r"))
	ca.decimals = uint8(bs[asset_decimals_position])
	ca.totalSupply = binary.LittleEndian.Uint64(bs[asset_total_supply_position:asset_price_position])
	ca.price = binary.LittleEndian.Uint64(bs[asset_price_position:asset_blocks_position])
	ca.blocks = binary.LittleEndian.Uint32(bs[asset_blocks_position:create_asset_from_position])
	ca.from = string(bs[create_asset_from_position:create_asset_nonce_position])
	ca.nonce = binary.LittleEndian.Uint32(bs[create_asset_nonce_position:create_asset_fee_position])
	ca.fee = binary.LittleEndian.Uint64(bs[create_asset_fee_position:create_asset_signer_position])
	ca.signer = decodeSigner(bs[create_asset_signer_position:])


	return ca
}

func (ca *createAsset_T) validate(fromP2p bool) error {
	txIdsMutex.Lock()
	defer txIdsMutex.Unlock()

	matched, err := regexp.MatchString("^[a-zA-Z0-9]{5,10}$", ca.name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("The length of `name` must between 5 to 10, and the characters must be littles or numbers")
	}

	if ca.name == "dsysb" || ca.name == "DSYSB" {
		return errors.New("`" + ca.name + "` has been kept")
	}

	matched, err = regexp.MatchString("^[a-zA-Z0-9]{3,5}$", ca.symbol)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("The length of `symbol` must between 3 to 5, and the characters must be littles or numbers")
	}

	if ca.symbol == "dsb" || ca.symbol == "DSB" {
		return errors.New("`" + ca.symbol + "` has been kept")
	}

	s := hex.EncodeToString(ca.signer.signature[:])
	if s == "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" {
		return errors.New("Unsigned transaction")
	}

	poolMutex.Lock()
	defer poolMutex.Unlock()

	// replay attack
	txIdH := ca.hash()
	txId := hex.EncodeToString(txIdH[:])
	for _, id := range txIds {
		if txId == id {
			if fromP2p {
				deleteFromTransactionPool(txId)
				return nil
			}
			return errors.New("Replay attack: txid: " + txId)
		}
	}

	var nonce uint32
	state, err := getState()
	if err != nil {
		return err
	}

	account, ok := state.accounts[ca.from]
	if ok {
		nonce = account.nonce
	}

	fmt.Println("nonce:", ca.nonce, nonce)
	if ca.nonce - nonce != 1 {
		return errors.New("The nonces are not match")
	}

	ok = ca.verifySign()
	if !ok {
		return errors.New("Invalid signature")
	}

	txIds = append(txIds, txId)

	return nil
}

func (ca *createAsset_T) verifySign() bool {
	publicKey := ecdsa.PublicKey{elliptic.P256(), ca.signer.x, ca.signer.y}
	txid := ca.hash()
	return ecdsa.Verify(&publicKey, txid[:], big.NewInt(0).SetBytes(ca.signer.signature[:32]), big.NewInt(0).SetBytes(ca.signer.signature[32:]))
}

func (ca *createAsset_T) count(state *state_T, coinbase *coinbase_T, index int) error {
	asset := &asset_T {
		ca.name,
		ca.symbol,
		ca.decimals,
		ca.totalSupply,
		ca.price,
		ca.blocks,
		ca.blocks,
	}

	assetIdB := asset.hash()
	assetId := hex.EncodeToString(assetIdB[:])
	_, ok := state.assets[assetId]
	if ok {
		return errors.New("Asset is already in")
	}

	account, ok := state.accounts[ca.from]
	if !ok {
		return errors.New("CA from is empty address")
	}

	holdAmount := ca.price * uint64(ca.blocks)
	totalSpend := holdAmount + ca.fee

	if account.balance < totalSpend {
		return errors.New("not enough minerals")
	}

	state.assets[assetId] = asset

	account.balance -= totalSpend
	coinbase.amount += ca.fee
	account.assets[assetId] = ca.totalSupply
	account.nonce = ca.nonce

	return nil
}

func (ca *createAsset_T) String() string {
	return fmt.Sprintf(
		"\tname: %s\n" +
		"\tsymbol: %s\n" +
		"\tdecimals: %d\n" +
		"\ttotol supply: %d\n" +
		"\tprice: %d\n" +
		"\tblocks: %d\n" +
		"\tfrom: %s\n" +
		"\tnonce: %d\n" +
		"\tfee: %d\n" +
		"%s", ca.name, ca.symbol, ca.decimals, ca.totalSupply, ca.price, ca.blocks, ca.from, ca.nonce, ca.fee, ca.signer)
}
