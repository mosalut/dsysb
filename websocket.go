// dsysb

package main

import (
	"net/http"

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
	address string // miner address
	head *blockHead_T
	poolCache *poolCache_T
}

func (wsAddBlockData *wsAddBlockData_T) encode() []byte {
	bs := append(wsAddBlockData.head.encode(), wsAddBlockData.poolCache.encode()...)

	return bs
}

func decodeWsAddBlockData(bs []byte) *wsAddBlockData_T {
	address := string(bs[:34])
	blockHead := decodeBlockHead(bs[34:bh_length + 34])
	poolCache := decodePoolCache(bs[34 + bh_length:])
	wsAddBlockData := &wsAddBlockData_T{
		address,
		blockHead,
		poolCache,
	}

	return wsAddBlockData
}

var upgrader = websocket.Upgrader{}

var minerConns = make(map[*websocket.Conn]interface{})
func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		print(log_error, "Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	minerConns[conn] = nil

	// The event loop
	for {
		data := socketData_T{}
		err := conn.ReadJSON(&data)
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				conn.Close()
				delete(minerConns, conn)
				print(log_warning, "close normal closure")
				break
			}
			print(log_error, err)
			break
		}

		switch data.Event {
		case WS_UPDATE:
			print(log_info, "update")

			address := string(data.Body)

			cache := poolToCache()
			coinbase := &coinbase_T {
				address,
				5e10,
			}
			cache.transactions = append([]transaction_I{ coinbase }, cache.transactions...)
			cache.count()
			bs := cache.encode()

			socketData := socketData_T { WS_UPDATE, bs }
			err = conn.WriteJSON(socketData)
			if err != nil {
				print(log_error, err)
				return
			}

			print(log_info, "ws_update sended")
		case WS_ADD_BLOCK:
			print(log_info, "new block")

			wsAddBlockData := decodeWsAddBlockData(data.Body)

			blockBody := &blockBody_T { wsAddBlockData.poolCache.transactions }
			block := &block_T { wsAddBlockData.head, blockBody }
			// TODO add block validation

			batch := &leveldb.Batch{}
			batch.Put([]byte("state"), wsAddBlockData.poolCache.state.encode())
			batch.Put(block.Head.Hash[32:], block.encode())

			if len(transactionPool) <= 511 {
				transactionPool = make([]transaction_I, 0, 511)
			} else {
				transactionPool = transactionPool[511:]
			}

			err = chainDB.Write(batch, nil)
			if err != nil {
				print(log_error, err)
				return
			}

			cache := poolToCache()
			coinbase := &coinbase_T {
				wsAddBlockData.address,
				5e10,
			}
			cache.transactions = append([]transaction_I{ coinbase }, cache.transactions...)
			cache.count()
			bs := cache.encode()

			socketData := socketData_T { WS_ADD_BLOCK, bs }

			for conn, _ := range minerConns {
				err = conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				print(log_info, conn.RemoteAddr(), "ws_state sended")
			}
		}
	}
}
