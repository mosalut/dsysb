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
	"errors"
	"fmt"
)

const (
	create_length = 203
	create_name_position = 1
	create_symbol_position = 11
	create_decimals_position = 16
	create_total_supply_position = 17
	create_price_position = 25
	create_blocks_position = 29
	create_from_position = 33
	create_nonce_position = 67
	create_byte_price_position = 71
	create_signer_position = 75
)

type createAsset_T struct {
	name string
	symbol string
	decimals uint8
	totalSupply uint64
	price uint32
	blocks uint32
	from string
	nonce uint32
	bytePrice uint32
	signer *signer_T
}

func (tx *createAsset_T) hash() [32]byte {
	bs := make([]byte, create_length, create_length)
	bs[0] = type_create
	copy(bs[create_name_position:create_symbol_position], []byte(tx.name))
	copy(bs[create_symbol_position:create_decimals_position], []byte(tx.symbol))
	bs[create_decimals_position] = byte(tx.decimals)
	binary.LittleEndian.PutUint64(bs[create_total_supply_position:create_price_position], tx.totalSupply)
	binary.LittleEndian.PutUint32(bs[create_price_position:create_blocks_position], tx.price)
	binary.LittleEndian.PutUint32(bs[create_blocks_position:create_from_position], tx.blocks)
	copy(bs[create_from_position:create_nonce_position], []byte(tx.from))
	binary.LittleEndian.PutUint32(bs[create_nonce_position:create_byte_price_position], tx.nonce)
	binary.LittleEndian.PutUint32(bs[create_byte_price_position:create_signer_position], tx.bytePrice)

	return sha256.Sum256(bs)
}

func (ca *createAsset_T) encode() []byte {
	bs := make([]byte, create_length, create_length)
	bs[0] = type_create
	copy(bs[create_name_position:create_symbol_position], []byte(ca.name))
	copy(bs[create_symbol_position:create_decimals_position], []byte(ca.symbol))
	bs[create_decimals_position] = byte(ca.decimals)
	binary.LittleEndian.PutUint64(bs[create_total_supply_position:create_price_position], ca.totalSupply)
	binary.LittleEndian.PutUint32(bs[create_price_position:create_blocks_position], ca.price)
	binary.LittleEndian.PutUint32(bs[create_blocks_position:create_from_position], ca.blocks)
	copy(bs[create_from_position:create_nonce_position], []byte(ca.from))
	binary.LittleEndian.PutUint32(bs[create_nonce_position:create_byte_price_position], ca.nonce)
	binary.LittleEndian.PutUint32(bs[create_byte_price_position:create_signer_position], ca.bytePrice)
	copy(bs[create_signer_position:], ca.signer.encode())

	return bs
}

func (tx *createAsset_T) encodeForPool() []byte {
	length := create_length + 2
	bs := make([]byte, length, length)
	binary.LittleEndian.PutUint16(bs[:2], create_length)
	copy(bs[2:], tx.encode())

	return bs
}

func decodeCreateAsset(bs []byte) *createAsset_T {
	ca := &createAsset_T{}

	ca.name = string(bytes.Trim(bs[create_name_position:create_symbol_position], "\x00 \t\n\r"))
	ca.symbol = string(bytes.Trim(bs[create_symbol_position:create_decimals_position], "\x00 \t\n\r"))
	ca.decimals = uint8(bs[create_decimals_position])
	ca.totalSupply = binary.LittleEndian.Uint64(bs[create_total_supply_position:create_price_position])
	ca.price = binary.LittleEndian.Uint32(bs[create_price_position:create_blocks_position])
	ca.blocks = binary.LittleEndian.Uint32(bs[create_blocks_position:create_from_position])
	ca.from = string(bs[create_from_position:create_nonce_position])
	ca.nonce = binary.LittleEndian.Uint32(bs[create_nonce_position:create_byte_price_position])
	ca.bytePrice = binary.LittleEndian.Uint32(bs[create_byte_price_position:create_signer_position])
	ca.signer = decodeSigner(bs[create_signer_position:])


	return ca
}

func (tx *createAsset_T) length() int {
	return create_length
}

func (tx *createAsset_T) fee() uint64 {
	return create_length * uint64(tx.bytePrice)
}

func (ca *createAsset_T) validate(head *blockHead_T, fromP2p bool) error {
	matched, err := regexp.MatchString("^[a-zA-Z0-9]{5,10}$", ca.name)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("The length of `name` must between 5 to 10, and the characters must be letters or numbers")
	}

	if ca.name == "dsysb" || ca.name == "DSYSB" {
		return errors.New("`" + ca.name + "` has been kept")
	}

	matched, err = regexp.MatchString("^[a-zA-Z0-9]{3,5}$", ca.symbol)
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("The length of `symbol` must between 3 to 5, and the characters must be letters or numbers")
	}

	if ca.symbol == "dsb" || ca.symbol == "DSB" {
		return errors.New("`" + ca.symbol + "` has been kept")
	}

	if ca.price == 0 {
		return errors.New("Asset's price must > 0")
	}

	/*
	if ca.blocks < 10000 {
		return errors.New("Asset's blocks must >= 10000")
	}
	*/

	if ca.bytePrice == 0 {
		return errors.New("Disallow zero byte price")
	}

	s := hex.EncodeToString(ca.signer.signature[:])
	if s == "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000" {
		return errors.New("Unsigned transaction")
	}

	// replay attack
	txIdH := ca.hash()
	txId := hex.EncodeToString(txIdH[:])
	for k, tx := range transactionPool {
		h := tx.hash()
		if txId == hex.EncodeToString(h[:]) {
			if fromP2p {
			//	deleteFromTransactionPool(txId)
				poolMutex.Lock()
				if len(transactionPool) - 1 == k {
					transactionPool = transactionPool[:k]
				} else {
					transactionPool = append(transactionPool[:k], transactionPool[k + 1:]...)
				}
				poolMutex.Unlock()
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
		return errNonceExpired
	}

	ok = ca.verifySign()
	if !ok {
		return errors.New("Invalid signature")
	}

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

	holdAmount := uint64(ca.price) * uint64(ca.blocks)
	totalSpend := holdAmount + ca.fee()

	if account.balance < totalSpend {
		return errors.New("not enough minerals")
	}

	state.assets[assetId] = asset

	account.balance -= totalSpend
	coinbase.amount += ca.fee()
	account.assets[assetId] = ca.totalSupply
	account.nonce = ca.nonce

	return nil
}

func (ca *createAsset_T) getBytePrice() uint32 {
	return ca.bytePrice
}

func (ca *createAsset_T) Map() map[string]interface{} {
	txM := make(map[string]interface{})
	h := ca.hash()
	txM["txid"] = hex.EncodeToString(h[:])
	txM["type"] = type_create
	txM["name"] = ca.name
	txM["symbol"] = ca.symbol
	txM["decimals"] = ca.decimals
	txM["totalSupply"] = ca.totalSupply
	txM["price"] = ca.price
	txM["blocks"] = ca.blocks
	txM["from"] = ca.from
	txM["nonce"] = ca.nonce
	txM["bytePrice"] = ca.bytePrice
	txM["signature"] = hex.EncodeToString(ca.signer.signature[:])

	return txM
}

func (ca *createAsset_T) String() string {
	return fmt.Sprintf(
		"txid:\t%064x\n" +
		"\tname: %s\n" +
		"\tsymbol: %s\n" +
		"\tdecimals: %d\n" +
		"\ttotol supply: %d\n" +
		"\tprice: %d\n" +
		"\tblocks: %d\n" +
		"\tfrom: %s\n" +
		"\tnonce: %d\n" +
		"\tbyte price: %d\n" +
		"\tfee: %d\n" +
		"%s", ca.hash(), ca.name, ca.symbol, ca.decimals, ca.totalSupply, ca.price, ca.blocks, ca.from, ca.nonce, ca.bytePrice, ca.fee(), ca.signer)
}
