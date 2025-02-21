// dsysb

package main

import (
	"math/big"
	"encoding/json"
	"crypto/elliptic"
	"io"
	"net"
	"net/http"
	"fmt"
	"errors"
)

const (
	type_coinbase = iota
	type_create
	type_transfer
	type_exchange
)

type publicKey_T struct {
	Curve *elliptic.CurveParams `json:"curve"`
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
}

type signer_T struct {
	PublicKey *publicKey_T `json:"PublicKey"`
	Signature [64]byte `json:"signature"`
}

func (signer *signer_T) String() string {
	return fmt.Sprintf(
		"public key:\t%x%x\n" +
		"signature:\t%x", signer.PublicKey.X.Bytes(), signer.PublicKey.Y.Bytes(), signer.Signature)
}

type transaction_T struct {
	Txid [32]byte `json:"txid"`
	Type uint8 `json:"type"`
	Data []byte `json:"data"`
}

func decodeRawTransaction(rawTransaction []byte) (*transaction_T, error) {
	transaction := transaction_T{}
	err := json.Unmarshal(rawTransaction, &transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func transactionValidate(transaction *transaction_T) error {
	switch transaction.Type {
	case type_coinbase:
		return errors.New("illage type")
	case type_create:
		ca, err := decodeCreateAsset(transaction.Data)
		if err != nil {
			return err
		}

		poolMutex.Lock()
		defer poolMutex.Unlock()
		for _, signature := range signatures {
			s := fmt.Sprintf("%0128x", ca.Signer.Signature)
			if s == signature {
				return errors.New(fmt.Sprintf("%064x", transaction.Txid) + " replay: " + s)
			}
			signatures = append(signatures, s)
		}

		var nonce uint32
		state := getState()
		account, ok := state.Accounts[ca.From]
		if ok {
			nonce = account.Nonce
		}

		if ca.Nonce - nonce != 1 {
			return errors.New("The nonces are not match")
		}
		return nil
	case type_transfer:
		transfer, err := decodeTransfer(transaction.Data)
		if err != nil {
			return err
		}

		if transfer.From == transfer.To {
			return errors.New("Transfer to self is not allowed")
		}

		for _, signature := range signatures {
			s := fmt.Sprintf("%0128x", transfer.Signer.Signature)
			if s == signature {
				return errors.New(fmt.Sprintf("%064x", transaction.Txid) + " replay: " + s)
			}
			signatures = append(signatures, s)
		}

		state := getState()
		assetId := fmt.Sprintf("%064x", transfer.AssetId)

		poolMutex.Lock()
		defer poolMutex.Unlock()

		if assetId != dsysbId {
			_, ok := state.Assets[assetId]
			if !ok {
				print(log_error, "There's not the asset id: " + assetId)
				return errors.New("There's not the asset id: " + assetId)
			}

		}

		var nonce uint32
		account, ok := state.Accounts[transfer.From]
		if !ok {
			return errors.New("There's not the account id")
		}

		nonce = account.Nonce
		fmt.Println("nonces:", transfer.Nonce, nonce)
		if transfer.Nonce - nonce != 1 {
			return errors.New("The nonces are not match")
		}

		return nil
	case type_exchange:
		exchange, err := decodeExchange(transaction.Data)
		if err != nil {
			return err
		}

		if exchange[0].From != exchange[1].To || exchange[0].To != exchange[1].From {
			return errors.New("Exchange address not match")
		}

		state := getState()
		poolMutex.Lock()
		defer poolMutex.Unlock()
		for _, transfer := range exchange {
			if transfer.From == transfer.To {
				return errors.New("Exchange to self is not allowed")
			}

			assetId := fmt.Sprintf("%064x", transfer.AssetId)

			if assetId != dsysbId {
				_, ok := state.Assets[assetId]
				if !ok {
					return errors.New("There's not the asset id: " + assetId)
				}
			}

			// proccess replay attack
			for _, signature := range signatures {
				s := fmt.Sprintf("%0128x", transfer.Signer.Signature)
				if s == signature {
					return errors.New(fmt.Sprintf("%064x", transaction.Txid) + " replay: " + s)
				}
				signatures = append(signatures, s)
			}

			var nonce uint32
			account, ok := state.Accounts[transfer.From]
			if !ok {
				return errors.New("There's not the account id")
			}

			nonce = account.Nonce
		//	fmt.Println("nonces:", transfer.Nonce, nonce)
			if transfer.Nonce - nonce != 1 {
				return errors.New("The nonces are not match")
			}
		}

		return nil
	}

	return nil
}

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

	transaction, err := decodeRawTransaction(params.Data)
	if err != nil {
		return err
	}

	err = transactionValidate(transaction)
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

	bs, err := json.Marshal(transactionPool)
	if err != nil {
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	writeResult(w, responseResult_T{true, "ok", bs})
}
