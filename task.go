// dsysb

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"errors"
)

type task_T struct {
	address string
	instructs []uint8
	vData []byte
}

func (task *task_T) encode() []byte {
	leng0 := len(task.instructs)
	leng1 := len(task.vData)
	leng := leng0 + leng1
	length := leng + 36 // 36 = address:34 + leng0:2
	bs := make([]byte, length, length)
	copy(bs[:34], []byte(task.address))
	binary.LittleEndian.PutUint16(bs[34:36], uint16(leng0))
	copy(bs[36:36 + leng0], []byte(task.instructs))
	copy(bs[36 + leng0:], task.vData)

	return bs
}

func (task *task_T) hash() [32]byte {
	leng := len(task.instructs)
	length := leng + 34
	bs := make([]byte, length, length)
	copy(bs[:34], []byte(task.address))
	copy(bs[34:34 + leng], []byte(task.instructs))

	return sha256.Sum256(bs)
}

func decodeTask(bs []byte) *task_T {
	task := &task_T{}
	task.address = string(bs[:34])
	leng0 := int(binary.LittleEndian.Uint16(bs[34:36]))
	task.instructs = bs[36:36 + leng0]
	task.vData = bs[36 + leng0:]

	return task
}

func (task *task_T) deploy() string {
	h := task.hash()
	key := hex.EncodeToString(h[:])
//	tasks = append(tasks, task)
	return key
}

func (task *task_T) excute(state *state_T) error {
	// var ip int for instructs
	/*
	defer func() {
		if r := recover(); r != nil {
			h := task.hash()
			print(log_warning, "task excute failed:", hex.EncodeToString(h[:]))
			print(r)
		}
	}()
	*/

	d := make([]byte, len(task.vData), len(task.vData))
	copy(d, task.vData)

	for ip := 0; ip < len(task.instructs); {
		switch task.instructs[ip] {
		case ins_movsb:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.movsb(p0, p1, p2)
		case ins_mov8:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mov8(p0, p1)
		case ins_mov16:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mov16(p0, p1)
		case ins_mov32:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mov32(p0, p1)
		case ins_mov64:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mov64(p0, p1)
		case ins_add8:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.add8(p0, p1, p2)
		case ins_add16:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.add16(p0, p1, p2)
		case ins_add32:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.add32(p0, p1, p2)
		case ins_add64:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.add64(p0, p1, p2)
		case ins_add8u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.add8u(p0, p1, p2)
		case ins_add16u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.add16u(p0, p1, p2)
		case ins_add32u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.add32u(p0, p1, p2)
		case ins_add64u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.add64u(p0, p1, p2)
		case ins_sub8:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub8(p0, p1, p2)
		case ins_sub16:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub16(p0, p1, p2)
		case ins_sub32:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub32(p0, p1, p2)
		case ins_sub64:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub64(p0, p1, p2)
		case ins_sub8u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub8u(p0, p1, p2)
		case ins_sub16u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub16u(p0, p1, p2)
		case ins_sub32u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub32u(p0, p1, p2)
		case ins_sub64u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.sub64u(p0, p1, p2)
		case ins_mul8:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul8(p0, p1, p2)
		case ins_mul16:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul16(p0, p1, p2)
		case ins_mul32:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul32(p0, p1, p2)
		case ins_mul64:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul64(p0, p1, p2)
		case ins_mul8u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul8u(p0, p1, p2)
		case ins_mul16u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul16u(p0, p1, p2)
		case ins_mul32u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul32u(p0, p1, p2)
		case ins_mul64u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.mul64u(p0, p1, p2)
		case ins_quo8:
			ip++
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
			ip++
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
			ip++
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
			ip++
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
			ip++
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
			ip++
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
			ip++
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
			ip++
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
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc8(p0)
		case ins_inc16:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc16(p0)
		case ins_inc32:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc32(p0)
		case ins_inc64:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc64(p0)
		case ins_inc8u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc8u(p0)
		case ins_inc16u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc16u(p0)
		case ins_inc32u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc32u(p0)
		case ins_inc64u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.inc64u(p0)
		case ins_dec8:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec8(p0)
		case ins_dec16:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec16(p0)
		case ins_dec32:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec32(p0)
		case ins_dec64:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec64(p0)
		case ins_dec8u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec8u(p0)
		case ins_dec16u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec16u(p0)
		case ins_dec32u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec32u(p0)
		case ins_dec64u:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			task.dec64u(p0)
		case ins_getIndex:
			ip++
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err := task.getIndex(p0)
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
