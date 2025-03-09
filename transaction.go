// dsysb

package main

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"fmt"
)

const (
	type_coinbase = iota
	type_create
	type_transfer
	type_exchange
)

type transaction_I interface {
	hash() [32]byte
	getType() uint8
	encode() []byte
	validate() error
	verifySign() bool
	count(*poolCache_T, int)
	String() string
}

func decodeRawTransaction(bs []byte) transaction_I {
	length := len(bs)

	var tx transaction_I
	switch length {
	case coinbase_length:
		tx = decodeCoinbase(bs)
	case create_asset_length:
		tx = decodeCreateAsset(bs)
	case transfer_length:
		tx = decodeTransfer(bs)
	case exchange_length:
		tx = decodeExchange(bs)
	default:
		print(log_error, "Wrong type")
	}

	return tx
}

/*
func transactionValidate(rawTransaction []byte) error {
	length := len(rawTransaction)

	switch length {
	case coinbase_length:
		return errors.New("illage type")
	case create_asset_length:
		ca := decodeCreateAsset(rawTransaction)

		poolMutex.Lock()
		defer poolMutex.Unlock()

		// replay attack
		for _, signature := range signatures {
			s := fmt.Sprintf("%0128x", ca.signer.signature)
			if s == signature {
				return errors.New(fmt.Sprintf("%064x", ca.hash()) + " replay: " + s)
			}
			signatures = append(signatures, s)
		}

		var nonce uint32
		state := getState()
		account, ok := state.accounts[ca.from]
		if ok {
			nonce = account.nonce
		}

		if ca.nonce - nonce != 1 {
			return errors.New("The nonces are not match")
		}
	case transfer_length:
		transfer := decodeTransfer(rawTransaction)

		if transfer.from == transfer.to {
			return errors.New("Transfer to self is not allowed")
		}

		for _, signature := range signatures {
			s := fmt.Sprintf("%0128x", transfer.signer.signature)
			if s == signature {
				return errors.New(fmt.Sprintf("%064x", transfer.hash()) + " replay: " + s)
			}
			signatures = append(signatures, s)
		}

		state := getState()
		assetId := fmt.Sprintf("%064x", transfer.assetId)

		poolMutex.Lock()
		defer poolMutex.Unlock()

		if assetId != dsysbId {
			_, ok := state.assets[assetId]
			if !ok {
				print(log_error, "There's not the asset id: " + assetId)
				return errors.New("There's not the asset id: " + assetId)
			}
		}

		var nonce uint32
		account, ok := state.accounts[transfer.from]
		if !ok {
			return errors.New("There's not the account id")
		}

		nonce = account.nonce
		fmt.Println("nonces:", transfer.nonce, nonce)
		if transfer.nonce - nonce != 1 {
			return errors.New("The nonces are not match")
		}
	case exchange_length:
		exchange := decodeExchange(rawTransaction)

		if exchange[0].from != exchange[1].to || exchange[0].to != exchange[1].from {
			return errors.New("Exchange address not match")
		}

		state := getState()
		poolMutex.Lock()
		defer poolMutex.Unlock()
		for _, transfer := range exchange {
			if transfer.from == transfer.to {
				return errors.New("Exchange to self is not allowed")
			}

			assetId := fmt.Sprintf("%064x", transfer.assetId)

			if assetId != dsysbId {
				_, ok := state.assets[assetId]
				if !ok {
					return errors.New("There's not the asset id: " + assetId)
				}
			}

			// proccess replay attack
			for _, signature := range signatures {
				s := fmt.Sprintf("%0128x", transfer.signer.signature)
				if s == signature {
					return errors.New(fmt.Sprintf("%064x", exchange.hash()) + " replay: " + s)
				}
				signatures = append(signatures, s)
			}

			var nonce uint32
			account, ok := state.accounts[transfer.from]
			if !ok {
				return errors.New("There's not the account id")
			}

			nonce = account.nonce
		//	fmt.Println("nonces:", transfer.Nonce, nonce)
			if transfer.nonce - nonce != 1 {
				return errors.New("The nonces are not match")
			}
		}
	default:
		return errors.New("Wrong length")

	}

	return nil
}
*/

func sendRawTransaction(body io.ReadCloser) error {
	defer body.Close()
	bs, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	params := &p2pParams_T{}
	err = json.Unmarshal(bs, &params)
	if err != nil {
		return err
	}
	params.Key = p2p_transport_sendrawtransaction_event

	bs, err = json.Marshal(params)
	if err != nil {
		return err
	}

	transaction := decodeRawTransaction(params.Data)

	/*
	err = transactionValidate(transaction)
	if err != nil {
		return err
	}
	*/

	err = transaction.validate()
	if err != nil {
		return err
	}

	// TODO  more validations
	poolMutex.Lock()
	transactionPool = append(transactionPool, transaction)
	poolMutex.Unlock()

	for k, _ := range seedAddrs {
		rAddr, err := net.ResolveUDPAddr("udp", k)
		if err != nil {
			print(log_error, err)
		}

		_, err = peer.Transport(rAddr, bs)
		if err != nil {
			print(log_error, err)
		}
	}

	return nil
}

func sendRawTransactionHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	err := sendRawTransaction(req.Body)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	writeResult(w, responseResult_T{true, "ok", nil})
}

func poolToCache() *poolCache_T {
	state := getState()

	fmt.Println("----------------------")
	fmt.Println(state.accounts)
	fmt.Println(state.accounts == nil)
	if len(transactionPool) <= 511 {
		return &poolCache_T {
			state,
			transactionPool,
		}
	}

	return &poolCache_T {
		state,
		transactionPool[:511],
	}
}

func txPool(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	fmt.Println(transactionPool)
	bs := transactionPool.encode()

	writeResult(w, responseResult_T{true, "ok", bs})
}
