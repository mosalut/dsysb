// dsysb

package main

import (
	"math/big"
	"fmt"
)

const (
	signer_length = 128
)

type signer_T struct {
	x *big.Int
	y *big.Int
	signature [64]byte
}

func (signer *signer_T) encode() []byte {
	bs := make([]byte, signer_length, signer_length)
	copy(bs[:32], signer.x.Bytes())
	copy(bs[32:64], signer.y.Bytes())
	copy(bs[64:], signer.signature[:])

	return bs
}

func decodeSigner(bs []byte) *signer_T {
	signer := &signer_T{}
	signer.x = big.NewInt(0)
	signer.x.SetBytes(bs[:32])
	signer.y = big.NewInt(0)
	signer.y.SetBytes(bs[32:64])
	signer.signature = [64]byte(bs[64:])

	return signer
}

func (signer *signer_T) String() string {
	return fmt.Sprintf(
		"\tpublic key:\t%x%x\n" +
		"\tsignature:\t%x", signer.x.Bytes(), signer.y.Bytes(), signer.signature)
}
