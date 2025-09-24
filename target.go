// dsysb

package main

import (
	"math/big"
	"encoding/binary"
)

var (
	stdBlockNum uint32
	stdBlockBatchSeconds int64
	difficult_1_target [4]byte
)

func initTargetValues() {
	if cmdFlag.networkID == 0x10 {
		// devnet
		stdBlockNum = 100 // for test faster
		stdBlockBatchSeconds = 60000 // 600 * 100 for dev faster
		difficult_1_target = [4]byte{ 0x1f, 0x00, 0xff, 0xff }
	} else if cmdFlag.networkID == 0 {
		// mainnet
		stdBlockNum = 2016
		stdBlockBatchSeconds = 1209600 // 600 * 2016
	//	difficult_1_target = [4]byte{ 0x1d, 0, 0xff, 0xff }
		difficult_1_target = [4]byte{ 0x1f, 0, 0xff, 0xff }

	} else {
		// testnet
		stdBlockNum = 1024
		stdBlockBatchSeconds = 614400 // 600 * 1024
	//	difficult_1_target = [4]byte{ 0x1d, 0, 0xff, 0xff }
		difficult_1_target = [4]byte{ 0x1f, 0, 0xff, 0xff }
	}
}

func adjustTarget(block *block_T) error {
	index := binary.LittleEndian.Uint32(block.head.hash[32:])
	if index % stdBlockNum != 1 || index < stdBlockNum {
		return nil
	}

	timestampNow := int64(binary.LittleEndian.Uint64(block.head.timestamp[:]))
	start := index - stdBlockNum

	startB := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(startB, start)

	startBlock, err := getBlock(startB)
	if err != nil {
		return err
	}

	timestampStart := int64(binary.LittleEndian.Uint64(startBlock.head.timestamp[:]))
	timestampDiff := timestampNow - timestampStart

	target := bitsToTarget(block.head.bits[:])
	x := big.NewInt(0)
	x = x.Mul(target, big.NewInt(timestampDiff))
	x.Div(x, big.NewInt(stdBlockBatchSeconds))

	if x.Cmp(bitsToTarget(difficult_1_target[:])) > 0 {
		block.head.bits = difficult_1_target
	} else {
		block.head.bits = targetToBits(x)
	}

	return nil
}

func bitsToTarget(bitsB []byte) *big.Int {
	target := big.NewInt(0).SetBytes(bitsB[1:])
	target.Lsh(target, 8 * uint(uint8(bitsB[0]) - 3))

	return target
}

func targetToBits(target *big.Int) [4]byte {
	var bits [4]byte
	targetB := target.Bytes()
	if uint8(targetB[0]) > 0x7f {
		bits[0] = byte(len(targetB) + 1)
		bits[1] = 0
		bits[2] = targetB[0]
		if len(targetB) == 1 {
			bits[3] = 0
		} else {
			bits[3] = targetB[1]
		}
	} else {
		bits[0] = byte(len(targetB))
		bits[1] = targetB[0]
		if len(targetB) == 1 {
			bits[2] = 0
			bits[3] = 0
		} else {
			bits[2] = targetB[1]
			if len(targetB) == 2 {
				bits[3] = 0
			} else {
				bits[3] = targetB[2]
			}
		}
	}

	return bits
}
