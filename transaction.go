package main

import (
	"math/big"
	"encoding/json"
	"encoding/hex"
//	"crypto/ecdsa"
	"crypto/elliptic"
	"io"
	"net"
	"net/http"
	"fmt"
)

type publicKey_T struct {
	Curve *elliptic.CurveParams `json:"curve"`
	X *big.Int `json:"x"`
	Y *big.Int `json:"y"`
}

type transaction_T struct {
	Txid []byte `json:"txid"`
	From string `json:"from"`
	To string `json:"to"`
	Script []byte `json:"script"`
	PublicKey *publicKey_T `json:"public_key"`
	Signature []byte `json:"signature"`
}

var transactionPool = make([]*transaction_T, 0, 32)

func (t transaction_T) String() string {
	return "txid:\t" + fmt.Sprintf("%x\n", t.Txid) +
	"from:\t" + t.From + "\n" +
	"to:\t" + t.To + "\n" +
	"script:\t" + hex.EncodeToString(t.Script) + "\n" +
	"signature:\t" + fmt.Sprintf("%x", t.Signature) + "\n"
}

func sendRawTransaction(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodPost:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	params := &p2pParams_T{}
	err = json.Unmarshal(body, &params)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}
	params.Key = p2p_transport_sendrawtransaction_event

	body, err = json.Marshal(params)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	dataStr := hex.EncodeToString(params.Data)

	transaction, err := decodeRawTransaction(dataStr)
	if err != nil {
		writeResult(w, responseResult_T{false, err.Error(), nil})
		return
	}

	transactionPool = append(transactionPool, transaction)
	fmt.Println(transaction)

	for k, _ := range seedAddrs {
		rAddr, err := net.ResolveUDPAddr("udp", k)
		if err != nil {
			writeResult(w, responseResult_T{false, err.Error(), nil})
			return
		}

		_, err = peer.Transport(rAddr, body)
		if err != nil {
			writeResult(w, responseResult_T{false, err.Error(), nil})
			return
		}
	}

	writeResult(w, responseResult_T{true, "ok", nil})
}

func decodeRawTransaction(s string) (*transaction_T, error) {
	rawTransaction, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	transaction := transaction_T{}
	err = json.Unmarshal(rawTransaction, &transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func getTransactionPool() []*transaction_T {
	return transactionPool
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
