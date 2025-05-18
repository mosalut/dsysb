// dsysb

package main

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
)

const (
	ins_movsb = iota
	ins_mov8
	ins_mov16
	ins_mov32
	ins_mov64
	ins_add8
	ins_add16
	ins_add32
	ins_add64
	ins_add8u
	ins_add16u
	ins_add32u
	ins_add64u
	ins_sub8
	ins_sub16
	ins_sub32
	ins_sub64
	ins_sub8u
	ins_sub16u
	ins_sub32u
	ins_sub64u
	ins_inc8
	ins_inc16
	ins_inc32
	ins_inc64
	ins_inc8u
	ins_inc16u
	ins_inc32u
	ins_inc64u
	ins_dec8
	ins_dec16
	ins_dec32
	ins_dec64
	ins_dec8u
	ins_dec16u
	ins_dec32u
	ins_dec64u
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
	return sha256.Sum256(task.encode())
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

func (task *task_T) excute() {
	// var ip int for instructs

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
			/*
		case ins_inc:
			ip++
		case ins_incU:
			ip++
		case ins_dec:
			ip++
		case ins_decU:
			ip++
			*/
		default:
		}
	}
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
