// dsysb

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"errors"
)

type reg_T struct {
	vBytes []byte
	vUint8 uint8
	vUint16 uint16
	vUint32 uint32
	vUint64 uint64
	vInt8 uint8
	vInt16 uint16
	vInt32 uint64
	vInt64 uint64

	vBool bool
}

type task_T struct {
	address string
	instructs []uint8
	vData []byte
	nonce uint32
	price uint32
	blocks uint32
	remain uint32
}

func (task *task_T) encode() []byte {
	leng0 := len(task.instructs)
	leng1 := len(task.vData)
	leng := leng0 + leng1
	length := leng + 52 // 52 = address:34 + leng0:2 + nonce:4 + price:4 + blocks:4 + remain:4
	bs := make([]byte, length, length)
	copy(bs[:34], []byte(task.address))
	binary.LittleEndian.PutUint16(bs[34:36], uint16(leng0))
	copy(bs[36:36 + leng0], []byte(task.instructs))
	copy(bs[36 + leng0:length - 16], task.vData)
	binary.LittleEndian.PutUint32(bs[length - 16:length - 12], task.nonce)
	binary.LittleEndian.PutUint32(bs[length - 12:length - 8], task.price)
	binary.LittleEndian.PutUint32(bs[length - 8:length - 4], task.blocks)
	binary.LittleEndian.PutUint32(bs[length - 4:], task.remain)

	return bs
}

func (task *task_T) hash() [32]byte {
	leng := len(task.instructs)
	length := leng + 46 // 46 = address:34 + leng:2 + nonce:4 + price:4 + blocks:4
	bs := make([]byte, length, length)
	copy(bs[:34], []byte(task.address))
	copy(bs[34:34 + leng], []byte(task.instructs))
	binary.LittleEndian.PutUint32(bs[length - 12:length - 8], task.nonce)
	binary.LittleEndian.PutUint32(bs[length - 8:length - 4], task.price)
	binary.LittleEndian.PutUint32(bs[length - 4:], task.blocks)

	return sha256.Sum256(bs)
}

func decodeTask(bs []byte) *task_T {
	length := len(bs)
	task := &task_T{}
	task.address = string(bs[:34])
	leng0 := int(binary.LittleEndian.Uint16(bs[34:36]))
	task.instructs = bs[36:36 + leng0]
	task.vData = bs[36 + leng0:length - 16]
	task.nonce = binary.LittleEndian.Uint32(bs[length - 16:length - 12])
	task.price = binary.LittleEndian.Uint32(bs[length - 12:length - 8])
	task.blocks = binary.LittleEndian.Uint32(bs[length - 8:length - 4])
	task.remain = binary.LittleEndian.Uint32(bs[length - 4:])

	return task
}

func (task *task_T) deploy() string {
	h := task.hash()
	key := hex.EncodeToString(h[:])
//	tasks = append(tasks, task) // for go testing
	return key
}

func (task *task_T) excute(state *state_T) error {
	// variable ip int for instructs

	d := make([]byte, len(task.vData), len(task.vData))
	copy(d, task.vData)

	reg := &reg_T{}

	for ip := 0; ip < len(task.instructs); {
		ipx := ip
		ip++
		switch task.instructs[ipx] {
		case ins_movsb:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // length
			ip += 2
			task.movsb(p0, p1, p2)
		case ins_mov8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			task.mov8(p0, p1)
		case ins_mov16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			task.mov16(p0, p1)
		case ins_mov32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			task.mov32(p0, p1)
		case ins_mov64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			task.mov64(p0, p1)
		case ins_add8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			task.add8(p0, p1, p2)
		case ins_add16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			task.add16(p0, p1, p2)
		case ins_add32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			task.add32(p0, p1, p2)
		case ins_add64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			task.add64(p0, p1, p2)
		case ins_add8u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			task.add8u(p0, p1, p2)
		case ins_add16u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			task.add16u(p0, p1, p2)
		case ins_add32u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			task.add32u(p0, p1, p2)
		case ins_add64u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			task.add64u(p0, p1, p2)
		case ins_sub8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub8(p0, p1, p2)
		case ins_sub16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub16(p0, p1, p2)
		case ins_sub32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub32(p0, p1, p2)
		case ins_sub64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub64(p0, p1, p2)
		case ins_sub8u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub8u(p0, p1, p2)
		case ins_sub16u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub16u(p0, p1, p2)
		case ins_sub32u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub32u(p0, p1, p2)
		case ins_sub64u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub64u(p0, p1, p2)
		case ins_mul8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul8(p0, p1, p2)
		case ins_mul16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul16(p0, p1, p2)
		case ins_mul32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul32(p0, p1, p2)
		case ins_mul64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul64(p0, p1, p2)
		case ins_mul8u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul8u(p0, p1, p2)
		case ins_mul16u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul16u(p0, p1, p2)
		case ins_mul32u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul32u(p0, p1, p2)
		case ins_mul64u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul64u(p0, p1, p2)
		case ins_quo8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.quo8(p0, p1, p2, p3)
		case ins_quo16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.quo16(p0, p1, p2, p3)
		case ins_quo32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.quo32(p0, p1, p2, p3)
		case ins_quo64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.quo64(p0, p1, p2, p3)
		case ins_quo8u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.quo8u(p0, p1, p2, p3)
		case ins_quo16u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.quo16u(p0, p1, p2, p3)
		case ins_quo32u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.quo32u(p0, p1, p2, p3)
		case ins_quo64u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.quo64u(p0, p1, p2, p3)
		case ins_inc8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc8(p0)
		case ins_inc16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc16(p0)
		case ins_inc32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc32(p0)
		case ins_inc64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc64(p0)
		case ins_inc8u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc8u(p0)
		case ins_inc16u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc16u(p0)
		case ins_inc32u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc32u(p0)
		case ins_inc64u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc64u(p0)
		case ins_dec8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec8(p0)
		case ins_dec16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec16(p0)
		case ins_dec32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec32(p0)
		case ins_dec64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec64(p0)
		case ins_dec8u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec8u(p0)
		case ins_dec16u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec16u(p0)
		case ins_dec32u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec32u(p0)
		case ins_dec64u:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec64u(p0)
		case ins_bytes:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.vBytes(reg, p0)
		case ins_uint8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination 
			ip += 2
			task.vUint8(reg, p0)
		case ins_uint16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.vUint16(reg, p0)
		case ins_uint32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.vUint32(reg, p0)
		case ins_uint64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.vUint64(reg, p0)
		case ins_int8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.vInt8(reg, p0)
		case ins_int16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.vInt16(reg, p0)
		case ins_int32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.vInt32(reg, p0)
		case ins_int64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.vInt64(reg, p0)
		case ins_eq:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.eq(reg, p0, p1, p2)
		case ins_gt:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.gt(reg, p0, p1, p2)
		case ins_lt:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.lt(reg, p0, p1, p2)
		case ins_gteq:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.gteq(reg, p0, p1, p2)
		case ins_lteq:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.lteq(reg, p0, p1, p2)
		case ins_eq_bytes:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.eqBytes(reg, p0, p1, p2)
		case ins_height:
			err := task.getIndex(reg)
			if err != nil {
				copy(task.vData, d)
				return err
			}
		default:
			if ip != 0 {
				copy(task.vData, d)
			}

			return errors.New("Invalid instruction")
		}
	}

	return nil
}

type taskPool_T []*task_T

func (pool taskPool_T) encode() []byte {
	bs := []byte{}
	for _, task := range pool {
		taskBytes := task.encode()
		leng := len(taskBytes)
		lengB := make([]byte, 2, 2)
		binary.LittleEndian.PutUint16(lengB, uint16(leng))
		bs = append(bs, lengB...)
		bs = append(bs, taskBytes...)
	}

	return bs
}

func decodeTaskPool(bs []byte) taskPool_T {
	pool := taskPool_T{}
	var currentStart int
	currentEnd := currentStart + 2
	length := len(bs)
	for currentEnd < length {
		taskBLength := int(binary.LittleEndian.Uint16(bs[currentStart:currentEnd]))
		currentStart = currentEnd
		currentEnd += taskBLength
		pool = append(pool, decodeTask(bs[currentStart:currentEnd]))
		currentStart = currentEnd
		currentEnd = currentEnd + 2
	}

	return pool
}

func tasksHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult(w, responseResult_T{false, "dsysb inner error", nil})
		return
	}

	writeResult(w, responseResult_T{true, "ok", state.tasks.encode()})
}

func taskHandler(w http.ResponseWriter, req *http.Request) {
	cors(w)

	switch req.Method {
	case http.MethodOptions:
		return
	case http.MethodGet:
	default:
		http.Error(w, API_NOT_FOUND, http.StatusNotFound)
		return
	}

	values := req.URL.Query()
	taskId := values.Get("id")

	state, err := getState()
	if err != nil {
		print(log_error, err)
		writeResult(w, responseResult_T{false, "dsysb inner error", nil})
		return
	}

	for _, task := range state.tasks {
		h := task.hash()
		tId := hex.EncodeToString(h[:])

		if tId == taskId {
			writeResult(w, responseResult_T{true, "ok", task.encode()})
			return
		}
	}

	writeResult(w, responseResult_T{false, "task " + taskId + " does not exist", nil})
}
