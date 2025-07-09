// dsysb

package main

import (
	"sort"
	"math/big"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"fmt"
)

const (
	asset_length = 36
	asset_symbol_position = 10
	asset_decimals_position = 15
	asset_total_supply_position = 16
	asset_price_position = 24
	asset_blocks_position = 28
	asset_remain_position = 32
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
	price uint32
	blocks uint32
	remain uint32
}

func (asset *asset_T) hash() [32]byte {
	bs := make([]byte, asset_length, asset_length)
	copy(bs[:asset_symbol_position], []byte(asset.name))
	copy(bs[asset_symbol_position:asset_decimals_position], []byte(asset.symbol))
	bs[asset_decimals_position] = byte(asset.decimals)
	binary.LittleEndian.PutUint64(bs[asset_total_supply_position:asset_price_position], asset.totalSupply)
	binary.LittleEndian.PutUint32(bs[asset_price_position:asset_blocks_position], asset.price)
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
	binary.LittleEndian.PutUint32(bs[asset_price_position:asset_blocks_position], asset.price)
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
	asset.price = binary.LittleEndian.Uint32(bs[asset_price_position:asset_blocks_position])
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
	length0 := len(pool)
	length := length0 * asset_length

	keys := make([]string, 0, length0)
	for k, _ := range pool {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		bsi := []byte(keys[i])
		bsj := []byte(keys[j])

		a := big.NewInt(0)
		a.SetBytes(bsi)

		b := big.NewInt(0)
		b.SetBytes(bsj)

		return a.Cmp(b) > 0
	})

	bs := make([]byte, 0, length)
	for _, key := range keys {
		bs = append(bs, pool[key].encode()...)
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
/*
type prolong_T struct {
	assetId [32]byte
	blocks [4]byte
	from [34]byte
}
*/

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
