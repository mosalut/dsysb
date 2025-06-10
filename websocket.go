// dsysb

package main

import (
	/* keepfunc */
	"encoding/hex"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/syndtr/goleveldb/leveldb" // keepfunc
)

const (
	WS_STATE = iota
	WS_UPDATE
	WS_MINED_BLOCK
	WS_ADD_BLOCK
	WS_ERR
)

type socketData_T struct {
	Event int `json:"event"`
	Body []byte `json:"body"`
}

type wsAddBlockData_T struct {
	head *blockHead_T
	poolCache *poolCache_T
}

func (wsAddBlockData *wsAddBlockData_T) encode() []byte {
	bs := append(wsAddBlockData.head.encode(), wsAddBlockData.poolCache.encode()...)

	return bs
}

func decodeWsAddBlockData(bs []byte) *wsAddBlockData_T {
	blockHead := decodeBlockHead(bs[:bh_length])
	poolCache := decodePoolCache(bs[bh_length:])
	wsAddBlockData := &wsAddBlockData_T{
		blockHead,
		poolCache,
	}

	return wsAddBlockData
}

var upgrader = websocket.Upgrader{}

var minerConns = make(map[*websocket.Conn]interface{})

/* keepfunc */
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

			cache, err := poolToCache()
			if err != nil {
				print(log_error, err)
				continue
			}
			cache.transactions = cache.transactions[:0]
			bs := cache.encode()

			socketData := socketData_T { WS_UPDATE, bs }
			err = conn.WriteJSON(socketData)
			if err != nil {
				print(log_error, err)
				continue
			}

			print(log_info, "ws_update sended")
		case WS_MINED_BLOCK:
			if blockchainSync.synchronizing {
				print(log_info, "new block, but synchronizing")
				continue
			}

			print(log_info, "new block")
			if len(data.Body) < 34 {
				print(log_error, "Data length wrong")
				socketData := socketData_T { WS_ERR, []byte("Data length wrong") }
				err := conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				continue
			}

			wsAddBlockData := decodeWsAddBlockData(data.Body)
			if wsAddBlockData == nil {
				socketData := socketData_T { WS_ERR, []byte("Wallet address format wrong") }
				err := conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				continue
			}

			blockBody := &blockBody_T { wsAddBlockData.poolCache.transactions }
			indexLength := bh_length + len(blockBody.encode())
			block := &block_T { wsAddBlockData.head, blockBody, wsAddBlockData.poolCache.state, uint32(indexLength)}

			// keepit
			//	TODO block validation
			lBlock, err := getHashBlock()
			if err != nil {
				print(log_error, err)
				continue
			}

			blockPrevHash := hex.EncodeToString(block.head.prevHash[:])
			lBlockHash := hex.EncodeToString(lBlock.head.hash[:])

			if blockPrevHash != lBlockHash {
				print(log_error, "The hash and prev hash are not match.")

				socketData, err := makeMinedBlockData()
				if err != nil {
					print(log_error, err)
					continue
				}

				err = conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				print(log_info, conn.RemoteAddr(), "ws_state sended")

				continue
			}

			if hex.EncodeToString(newMerkleTree(block.body.transactions).data[:]) != hex.EncodeToString(block.head.transactionRoot[:]) {
				print(log_warning, "WS: The transactionRoot and it's data are not match.")

				socketData, err := makeMinedBlockData()
				if err != nil {
					print(log_error, err)
					continue
				}

				err = conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				print(log_info, conn.RemoteAddr(), "ws_state sended")

				continue
			}

			h := block.state.hash()
			if hex.EncodeToString(h[:]) != hex.EncodeToString(block.head.stateRoot[:]) {
				print(log_warning, "WS: The stateRoot and it's data are not match.")

				socketData, err := makeMinedBlockData()
				if err != nil {
					print(log_error, err)
					continue
				}

				err = conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				print(log_info, conn.RemoteAddr(), "ws_state sended")

				continue
			}
			// ------------------------------------------

			batch := &leveldb.Batch{}
			batch.Put([]byte("index"), block.head.hash[32:])
			batch.Put(block.head.hash[32:], block.encode())

			err = chainDB.Write(batch, nil)
			if err != nil {
				print(log_error, err)
				continue
			}

			socketData, err := makeMinedBlockData()
			if err != nil {
				print(log_error, err)
				continue
			}

			for conn, _ := range minerConns {
				err = conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				print(log_info, conn.RemoteAddr(), "ws_state sended")
			}

			blockData := block.encode()
			broadcast(p2p_add_block_event, blockData)
		}
	}
}

func makeMinedBlockData() (*socketData_T, error) {
	cache, err := poolToCache()
	if err != nil {
		return nil, err
	}

	bs := cache.encode()

	if len(transactionPool) <= 511 {
		transactionPool = make([]transaction_I, 0, 511)
	} else {
		transactionPool = transactionPool[511:]
	}

	txIdsMutex.Lock()
	txIds = make([]string, 0, 511)
	txIdsMutex.Unlock()

	return &socketData_T { WS_MINED_BLOCK, bs }, nil
}
