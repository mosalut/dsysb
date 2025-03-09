// dsysb

package main

import (
	"log"
)

type blockchain_T []*blockHead_T

func (chain blockchain_T) encode() []byte {
	length := len(chain)
	lenH := length * bh_length
	bs := make([]byte, lenH, lenH)

	for k, head := range chain {
		copy(bs[k * bh_length:(k + 1) * bh_length], head.encode())
	}

	return bs
}

func decodeBlockchain(bs []byte) blockchain_T {
	length := len(bs)
	if length % bh_length != 0 {
		log.Fatal("Wrong length of block heads")
	}

	lenH := length / bh_length
	blockchain := make([]*blockHead_T, lenH, lenH)
	for i := 0; i < lenH; i++ {
		head := decodeBlockHead(bs[i * bh_length:(i + 1) * bh_length])
		blockchain[i] = head
	}

	return blockchain
}
