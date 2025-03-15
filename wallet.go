package main

import (
	"crypto/sha256"
	"crypto/ecdsa"
	"bytes"

	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcutil/base58"
	"github.com/syndtr/goleveldb/leveldb"
)

const walletVersion = byte(0x1e)
const addressChecksumLen = 4

var walletDB *leveldb.DB

type wallet_T struct {
	privateKey ecdsa.PrivateKey
}

// TODO
/*
func (w wallet_T) verify(, dataStr string) []byte {
	data, _ := hex.DecodeString(dataStr)
}
*/

func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		print(1, err)
		return nil
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

func validateAddress(address string) bool {
	pubKeyHash := base58.Decode(address)
	pubKeyHashLength := len(pubKeyHash)
	if pubKeyHashLength != 25 {
		return false
	}
	actualChecksum := pubKeyHash[pubKeyHashLength - addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : pubKeyHashLength - addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
