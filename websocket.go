// dsysb

package main

import (
	"encoding/json"
	"net/http"
//	"fmt"

	"github.com/gorilla/websocket"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	WS_STATE = iota
	WS_UPDATE
	WS_ADD_BLOCK
)

type socketData_T struct {
	Event int `json:"event"`
	Body []byte `json:"body"`
}

type wsAddBlockData_T struct {
	Head *blockHead_T `json:"head"`
	PoolCache *poolCache_T `json:"poolCache"`
}

func (wsAddBlockData *wsAddBlockData_T) encode() ([]byte, error) {
	bs, err := json.Marshal(wsAddBlockData)
	if err != nil {
		return nil, err
	}

	return bs, nil
}

func decodeWsAddBlockData(bs []byte) (*wsAddBlockData_T, error) {
	wsAddBlockData := &wsAddBlockData_T{}
	err := json.Unmarshal(bs, wsAddBlockData)
	if err != nil {
		return nil, err
	}

	return wsAddBlockData, err
}

var upgrader = websocket.Upgrader{}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		print(log_error, "Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	// The event loop
	for {
		data := socketData_T{}
		err := conn.ReadJSON(&data)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				conn.Close()
				print(log_warning, "close normal closure")
				break
			}
			print(log_error, err)
			break
		}

		switch data.Event {
		case WS_UPDATE:
			print(log_info, "update")

			cache := poolToCache()
			bs, err := cache.encode()
			if err != nil {
				print(log_error, err)
				return
			}

			socketData := socketData_T { WS_UPDATE, bs }
			err = conn.WriteJSON(socketData)
			if err != nil {
				print(log_error, err)
				return
			}

			print(log_info, "ws_update sended")
		case WS_ADD_BLOCK:
			print(log_info, "new block")

			wsAddBlockData, err := decodeWsAddBlockData(data.Body[:])
			if err != nil {
				print(log_error, err)
				return
			}

			blockBody := &blockBody_T { wsAddBlockData.PoolCache.Transactions }
			block := &block_T { wsAddBlockData.Head, blockBody }
			// TODO add block validation

			batch := &leveldb.Batch{}
			batch.Put([]byte("state"), wsAddBlockData.PoolCache.State.encode())
			batch.Put(block.Head.Hash[32:], block.encode())

			if len(transactionPool) <= 511 {
				transactionPool = make([]*transaction_T, 0, 511)
			} else {
				transactionPool = transactionPool[511:]
			}

			err = chainDB.Write(batch, nil)
			if err != nil {
				print(log_error, err)
				return
			}

			poolCache := poolToCache()
			bs, err := poolCache.encode()
			if err != nil {
				print(log_error, err)
				return
			}

			socketData := socketData_T { WS_ADD_BLOCK, bs }
			err = conn.WriteJSON(socketData)
			if err != nil {
				print(log_error, err)
				return
			}

			print(log_info, "ws_state sended")
		}
	}
}
