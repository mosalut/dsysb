// dsysb

package main

import (
	"bytes"
	"crypto/sha256"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"encoding/binary"
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

	create_asset_length = 202
	create_asset_from_position = 36
	create_asset_nonce_position = 70
	create_asset_signer_position = 74
)

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
		pool[fmt.Sprintf("%064x", asset.hash())] = asset
	}

	return pool
}

// `Blocks` represent how many Blocks the asset holden.
type prolong_T struct {
	assetId [32]byte
	blocks [4]byte
	from [34]byte
}

func listAssetsHandler(w http.ResponseWriter, req *http.Request) {
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

	fmt.Println(state.assets)

	writeResult(w, responseResult_T{true, "ok", state.assets.encode()})
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
	binary.LittleEndian.PutUint32(bs[create_asset_nonce_position:create_asset_signer_position], tx.nonce)

	return sha256.Sum256(bs)
}

func (ca *createAsset_T) getType() uint8 {
	return type_create
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
	binary.LittleEndian.PutUint32(bs[create_asset_nonce_position:create_asset_signer_position], ca.nonce)
	copy(bs[create_asset_signer_position:], ca.signer.encode())

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
	ca.nonce = binary.LittleEndian.Uint32(bs[create_asset_nonce_position:create_asset_signer_position])
	ca.signer = decodeSigner(bs[create_asset_signer_position:])


	return ca
}

func (ca *createAsset_T) validate(fromP2p bool) error {
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

	s := fmt.Sprintf("%0128x", ca.signer.signature)
	if s == "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" {
		return errors.New("Unsigned transaction")
	}

	poolMutex.Lock()
	defer poolMutex.Unlock()

	// replay attack
	for _, signature := range signatures {
		if s == signature {
			return errors.New("Replay attack: hash:" + fmt.Sprintf("%064x", ca.hash()) + " signature: " + s)
		}
	}
	signatures = append(signatures, s)

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

	return nil
}

func (ca *createAsset_T) verifySign() bool {
	publicKey := ecdsa.PublicKey{elliptic.P256(), ca.signer.x, ca.signer.y}
	fmt.Println(publicKey)
	txid := ca.hash()
	ok := ecdsa.Verify(&publicKey, txid[:], big.NewInt(0).SetBytes(ca.signer.signature[:32]), big.NewInt(0).SetBytes(ca.signer.signature[32:]))
	return ok
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
		"%s", ca.name, ca.symbol, ca.decimals, ca.totalSupply, ca.price, ca.blocks, ca.from, ca.nonce, ca.signer)
}
