// dsysb 

package main

import (
	"sort"
	"math/big"
	"crypto/sha256"
	"encoding/binary"
	"net/http"
	"fmt"
)

type state_T struct {
	assets assetPool_T
	accounts map[string]*account_T
	tasks []*task_T
}

func (state *state_T)encode() []byte {
	assetLength := len(state.assets) * asset_length

	var accountLength int
	for _, account := range state.accounts {
		accountLength += 14 + len(account.assets) * 40 + 34 // 14 = 8 + 4 + 2, 40 = 32 + 8
	}

	length := 4 + assetLength + accountLength
	bs := make([]byte, length, length)
	var start int

	end := assetLength
	copy(bs[start:end], state.assets.encode())

	acsLength := len(state.accounts)
	keys := make([]string, 0, acsLength)
	for k, _ := range state.accounts {
		keys = append(keys, k)
	}

	// Cause It can not ensure map's order in each reading
	sort.Slice(keys, func(i, j int) bool {
		bsi := []byte(keys[i])
		bsj := []byte(keys[j])

		a := big.NewInt(0)
		a.SetBytes(bsi)

		b := big.NewInt(0)
		b.SetBytes(bsj)

		return a.Cmp(b) > 0
	})

	for _, key := range keys {
		start = end
		end += 34
		copy(bs[start:end], key)

		accountBytes := state.accounts[key].encode()
		start = end
		end += len(accountBytes)
		copy(bs[start:end], accountBytes)

		start = end
		end += 2
		binary.LittleEndian.PutUint16(bs[start:end], uint16(len(state.accounts[key].assets)))
	}

	start = end

	fmt.Println("accountLength:", accountLength)
	binary.LittleEndian.PutUint32(bs[start:], uint32(accountLength))

	// encoding tasks
	var tasksBytesLength uint32
	for _, task := range state.tasks {
		taskBytes := task.encode()
		leng := len(taskBytes)
		tasksBytesLength += uint32(leng) + 2
		lengB := make([]byte, 2, 2)
		binary.LittleEndian.PutUint16(lengB, uint16(leng))
		bs = append(bs, lengB...)
		bs = append(bs, taskBytes...)
	}

	tasksBytesLengthB := make([]byte, 4, 4)
	binary.LittleEndian.PutUint32(tasksBytesLengthB, tasksBytesLength)
	bs = append(bs, tasksBytesLengthB...)

	return bs
}

func decodeState(bs []byte) *state_T {
	fmt.Println("bs length:", len(bs))
	var start, end int

	state := &state_T{}
	state.tasks = []*task_T{}

	start = len(bs) - 4
	tasksBytesLength := int(binary.LittleEndian.Uint32(bs[start:]))
	tasksStartPosition := len(bs) - tasksBytesLength - 4

	currentStart := tasksStartPosition
	currentEnd := currentStart + 2
	for currentEnd < start {
		taskBLength := int(binary.LittleEndian.Uint16(bs[currentStart:currentEnd]))
		state.tasks = append(state.tasks, decodeTask(bs[currentEnd:taskBLength]))
		currentStart = currentEnd + taskBLength
		currentEnd = currentStart + 2
	}

	start = tasksStartPosition - 4
	fmt.Println("start:", start)
	fmt.Println("tasksStartPosition:", tasksStartPosition)
	accountBytesLength := int(binary.LittleEndian.Uint32(bs[start:tasksStartPosition]))
	fmt.Println("accountBytesLength:", accountBytesLength)
	assetEndPosition := len(bs) - accountBytesLength - 4

	state.assets = decodeAssetPool(bs[0:assetEndPosition])

	state.accounts = make(map[string]*account_T)

	var assetsInAccount int

	for start > assetEndPosition {
		end = start
		start = end - 2

		assetsInAccount = int(binary.LittleEndian.Uint16(bs[start:end]))

		end = start
		start -= 12 + assetsInAccount * 40
		accountBytes := bs[start:end]

		end = start
		start -= 34
		state.accounts[string(bs[start:end])] = decodeAccount(accountBytes)
	}

	return state
}

func (state *state_T)hash() [32]byte {
	return sha256.Sum256(state.encode())
}

func (state *state_T)String() string {
	value := ("state:\n")
	value += "\tassets:\n"
	for _, asset := range state.assets {
		value += fmt.Sprintf("\t\t%v\n", asset)
	}

	value += "\taccounts:\n"
	for k, account := range state.accounts {
		value += fmt.Sprintf("\t\t%s: %v\n", k, account)
	}

	value += "\ttasks:\n"
	for _, task := range state.tasks {
		value += fmt.Sprintf("\t\t%s:\n", task.hash())
		value += fmt.Sprintf("\t\t%v:\n", task)
	}

	return value

}

var firstState = &state_T {
//	binary.LittleEndian.Uint32(difficult_1_target[:]),
	make(assetPool_T),
	make(map[string]*account_T),
	make([]*task_T, 0),
}

func getIndex() (uint32, error) {
	indexB, err := chainDB.Get([]byte("index"), nil)
	if err != nil {
		return 0, err
	}

	index := binary.LittleEndian.Uint32(indexB)

	return index, err
}

func getState() (*state_T, error) {
	indexB, err := chainDB.Get([]byte("index"), nil)
	if err != nil {
		return nil, err
	}

	/* keepit
//	indexB := []byte{0, 0, 0, 0}
	return firstState, nil
	*/

	block, err := getBlock(indexB)
	if err != nil {
		return nil, err
	}

	return block.state, nil
}

func stateHandler(w http.ResponseWriter, req *http.Request) {
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

	stateBytes := state.encode()

	writeResult(w, responseResult_T{true, "ok", stateBytes})
}
