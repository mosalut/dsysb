package main

import (
	"crypto/sha256"
)

type merkleNode_T struct {
	left *merkleNode_T
	right *merkleNode_T
	data [32]byte
}

func newMerkleTree(transactions []transaction_I) *merkleNode_T {
	hashes := make([][32]byte, len(transactions), len(transactions))
	for k, transaction := range transactions {
		hashes[k] = transaction.hash()
	}

	var level []*merkleNode_T
	if len(hashes) % 2 != 0 {
		level = make([]*merkleNode_T, len(hashes), len(hashes))
	} else {
		level = make([]*merkleNode_T, len(hashes) + 1, len(hashes) + 1)
	}

	for k, hash := range hashes {
		level[k] = &merkleNode_T{}
		level[k].left = nil
		level[k].right = nil
		level[k].data = hash
	}

	for {
		if len(level) == 1 {
			return level[0]
		}

		var newLevel []*merkleNode_T
		if len(level) % 2 != 0 {
			newLevel = make([]*merkleNode_T, len(hashes), len(hashes))
		} else {
			newLevel = make([]*merkleNode_T, len(hashes) + 1, len(hashes) + 1)
		}

		for k, _ := range newLevel {
			newLevel[k] = &merkleNode_T{}
			newLevel[k].left = level[k * 2]
			newLevel[k].right = level[k * 2 + 1]
			newLevel[k].data = sha256.Sum256(append(newLevel[k].left.data[:], newLevel[k].right.data[:]...))
		}

		level = newLevel
	}
}
