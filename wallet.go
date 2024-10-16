package main

import (
	"crypto/sha256"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	"bytes"

	"golang.org/x/crypto/ripemd160"
	"github.com/btcsuite/btcutil/base58"
)

const walletVersion = byte(0x01)
const addressChecksumLen = 4

type wallet_T struct {
	privateKey ecdsa.PrivateKey
	publicKey  []byte
}

func (w wallet_T) getAddress() []byte {
	pubKeyHash := HashPubKey(w.publicKey)

	versionedPayload := append([]byte{walletVersion}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := base58.Encode(fullPayload)

	return []byte(address)
}

func newWallet() *wallet_T {
	private, public := newKeyPair()
	wallet := wallet_T{*private, public}

	return &wallet
}

func newKeyPair() (*ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		print(1, err)
		return nil, nil
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return private, pubKey
}

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

type wallets_T struct {
	wallets map[string]*wallet_T
}

func validateAddress(address string) bool {
	pubKeyHash := base58.Decode(address)
	actualChecksum := pubKeyHash[len(pubKeyHash) - addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash) - addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
