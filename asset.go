// dsysb

package main

import (
	"crypto/sha256"
	"encoding/json"
	"net/http"
	"errors"
	"fmt"
)

// The `Name` is the asset Name.
// The `Symbol` is the asset Symbol.
// `TotalSupply` is just total supply.
// `Price` represent the Price of the asset that is holden by an block.
type asset_T struct {
	Name string `json:"name"`
	Symbol string `json:"symbol"`
	Decimals uint8 `json:"decimals"`
	TotalSupply uint64 `json:"totalSupply"`
	Price uint64 `json:"price"`
	Blocks uint32 `json:"blocks"`
	Height uint32 `json:"height"`
}

func (asset *asset_T) hash () ([32]byte, error) {
	bs, err := asset.encode()
	if err != nil {
		return [32]byte{}, err
	}
	h := sha256.Sum256(bs)

	return h, nil
}

func (asset *asset_T) encode () ([]byte, error) {
	bs, err := json.Marshal(asset)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func decodeAsset(bs []byte) (*asset_T, error) {
	asset := &asset_T{}
	err := json.Unmarshal(bs, asset)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

var assetPool = make([]*asset_T, 0, 500)

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

	state := getState()

	bs, err := json.Marshal(state.Assets)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	writeResult(w, responseResult_T{true, "ok", bs})
}

type createAsset_T struct {
	Name string `json:"name"`
	Symbol string `json:"symbol"`
	Decimals uint8 `json:"decimals"`
	TotalSupply uint64 `json:"totalSupply"`
	Price uint64 `json:"price"`
	Blocks uint32 `json:"blocks"`
	From string `json:"from"`
	Nonce uint32 `json:"nonce"`
	Signer *signer_T `json:"signer"`
}

func (ca *createAsset_T) encode() ([]byte, error) {
	bs, err := json.Marshal(ca)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func decodeCreateAsset(bs []byte) (*createAsset_T, error) {
	ca := &createAsset_T{}
	err := json.Unmarshal(bs, ca)
	if err != nil {
		return nil, err
	}

	return ca, nil
}

func (ca *createAsset_T) check() error {
	if len(ca.Name) > 10 || len(ca.Name) < 5 {
		return errors.New("The length of `name` must between 5 to 10")
	}

	if len(ca.Symbol) > 5 || len(ca.Symbol) < 3 {
		return errors.New("The length of `name` must between 3 to 5")
	}

	if ca.Decimals > 18 {
		return errors.New("`decimals` must <= 18")
	}

	if ca.Price == 0 {
		return errors.New("`price` must > 0")
	}

	if ca.Blocks < 10000 {
		return errors.New("`blocks` must >= 10000")
	}

	if !validateAddress(ca.From) {
		return errors.New("`from`: invalid address")
	}

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
		"%s", ca.Name, ca.Symbol, ca.Decimals, ca.TotalSupply, ca.Price, ca.Blocks, ca.From, ca.Nonce, ca.Signer)
}
