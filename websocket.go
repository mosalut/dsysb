// dsysb

package main

import (
	"encoding/binary"
	"encoding/hex"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/syndtr/goleveldb/leveldb" // keepfunc
)

const (
	WS_START = iota
	WS_MINED_BLOCK
	WS_PREPARED_BLOCK
	WS_ERR
)

const (
	WS_NOTICE_PING = iota
	WS_NOTICE_APPEND
	WS_NOTICE_TASK
	WS_NOTICE_ERR
)

type socketData_T struct {
	Event int `json:"event"`
	Body []byte `json:"body"`
}

type noticeData_T struct {
	Event int `json:"event"`
	Body interface{} `json:"body"`
}

var upgrader = websocket.Upgrader {
	CheckOrigin: func(r *http.Request) bool {
            return true
        },
}

var minerConns = make(map[*websocket.Conn]interface{})
var noticeConns = make(map[*websocket.Conn]interface{})

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
			if isWebsocketCloseError(err) {
				conn.Close()
				delete(minerConns, conn)
				print(log_warning, err, "close normal closure")
				break
			}
			print(log_error, err)
			break
		}

		switch data.Event {
		case WS_START:
			print(log_info, "start")

			if blockchainSync.synchronizing {
				err604 := minedError{"604", errSynchronizing.Error()}
				noticeErrorBroadcast(err604)
				print(log_info, err604)
				continue
			}

			address := string(data.Body)

			block, err := makeBlockForMine(address)
			if err != nil {
				err607 := minedError{"607", err.Error()}
				noticeErrorBroadcast(err607)
				print(log_error, err607)

				socketData := socketData_T { WS_ERR, []byte(err.Error()) }
				err0 := conn.WriteJSON(socketData)
				if err0 != nil {
					print(log_error, err0)
					continue
				}

				continue
			}
			bs := block.encode()

			socketData := socketData_T { WS_START, bs }
			err = conn.WriteJSON(socketData)
			if err != nil {
				print(log_error, err)
				continue
			}

			print(log_info, "ws_update sent")
		case WS_MINED_BLOCK:
			if blockchainSync.synchronizing {
				err604 := minedError{"604", errSynchronizing.Error()}
				noticeErrorBroadcast(err604)
				print(log_info, err604)
				continue
			}

			print(log_info, "new block")

			block, err := decodeBlock(data.Body)
			if err != nil {
				noticeErrorBroadcast(err)
				print(log_error, err)
				continue
			}

			// keepit
			// block validation
			lBlock, err := getHashBlock()
			if err != nil {
				err606 := minedError{"606", err.Error()}
				noticeErrorBroadcast(err606)
				print(log_error, err606)
				continue
			}

			blockPrevHash := hex.EncodeToString(block.head.prevHash[:])
			lBlockHash := hex.EncodeToString(lBlock.head.hash[:])

			if blockPrevHash != lBlockHash {
				err601 := minedError{"601", errPrevHashNotMatch.Error()}
				noticeErrorBroadcast(err601)
				print(log_error, err601)

				block, err := makeBlockForMine(block.body.transactions[0].(*coinbase_T).to)
				if err != nil {
					err607 := minedError{"607", err.Error()}
					noticeErrorBroadcast(err607)
					print(log_error, err607)

					socketData := socketData_T { WS_ERR, []byte(err.Error()) }
					err0 := conn.WriteJSON(socketData)
					if err0 != nil {
						print(log_error, err0)
						continue
					}

					continue
				}

				socketData := socketData_T { WS_PREPARED_BLOCK, block.encode() }
				err = conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				print(log_info, conn.RemoteAddr(), "ws_state sent")

				continue
			}

			if hex.EncodeToString(newMerkleTree(block.body.transactions).data[:]) != hex.EncodeToString(block.head.transactionRoot[:]) {
				err602 := minedError{"602", errTransactionRootNotMatch.Error()}
				noticeErrorBroadcast(err602)
				print(log_error, err602)

				block, err := makeBlockForMine(block.body.transactions[0].(*coinbase_T).to)
				if err != nil {
					err607 := minedError{"607", err.Error()}
					noticeErrorBroadcast(err607)
					print(log_error, err607)

					socketData := socketData_T { WS_ERR, []byte(err.Error()) }
					err0 := conn.WriteJSON(socketData)
					if err0 != nil {
						print(log_error, err0)
						continue
					}

					continue
				}

				socketData := socketData_T { WS_PREPARED_BLOCK, block.encode() }
				err = conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				print(log_info, conn.RemoteAddr(), "ws_state sent")

				continue
			}

			h := block.state.hash()
			if hex.EncodeToString(h[:]) != hex.EncodeToString(block.head.stateRoot[:]) {
				err603 := minedError{"603", errStateRootNotMatch.Error()}
				noticeErrorBroadcast(err603)
				print(log_error, err603)

				block, err := makeBlockForMine(block.body.transactions[0].(*coinbase_T).to)
				if err != nil {
					err607 := minedError{"607", err.Error()}
					noticeErrorBroadcast(err607)
					print(log_error, err607)

					socketData := socketData_T { WS_ERR, []byte(err.Error()) }
					err0 := conn.WriteJSON(socketData)
					if err0 != nil {
						print(log_error, err0)
						continue
					}

					continue
				}

				socketData := socketData_T { WS_PREPARED_BLOCK, block.encode() }
				err = conn.WriteJSON(socketData)
				if err != nil {
					print(log_error, conn.RemoteAddr, err)
					continue
				}
				print(log_info, conn.RemoteAddr(), "ws_state sent")

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

			block, err = makeBlockForMine(block.body.transactions[0].(*coinbase_T).to)
			if err != nil {
				err607 := minedError{"607", err.Error()}
				noticeErrorBroadcast(err607)
				print(log_error, err607)

				continue
			}

			socketData := socketData_T { WS_PREPARED_BLOCK, block.encode() }
			err = conn.WriteJSON(socketData)
			if err != nil {
				print(log_error, conn.RemoteAddr, err)
				continue
			}

			print(log_info, conn.RemoteAddr(), "ws_state sent")

			broadcast(p2p_add_block_event, data.Body)

			socketData = socketData_T { WS_MINED_BLOCK, nil }
			for c, _ := range minerConns {
				if c.RemoteAddr() == conn.RemoteAddr() {
					continue
				}

				err = c.WriteJSON(socketData)
				if err != nil {
					print(log_error, c.RemoteAddr(), err)
					continue
				}
				print(log_info, c.RemoteAddr(), "ws_state mined sent")
			}
		}
	}
}

// notice
func noticeSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		print(log_error, "Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	noticeConns[conn] = nil

	for {
		data := noticeData_T{}
		err := conn.ReadJSON(&data)
		if err != nil {
			if isWebsocketCloseError(err) {
				conn.Close()
				delete(noticeConns, conn)
				print(log_warning, err, "close normal closure")
				break
			}
			print(log_error, err)
			break
		}
	}
}

func noticeError(conn *websocket.Conn, msg string) {
	noticePush(conn, WS_NOTICE_ERR, msg)
}

func noticePush(conn *websocket.Conn, event int, data interface{}) error {
	noticeData := noticeData_T { event, data }

	err := conn.WriteJSON(noticeData)
	if err != nil {
		if isWebsocketCloseError(err) {
			conn.Close()
			delete(noticeConns, conn)
			print(log_warning, conn.RemoteAddr(), err, "close normal closure")
			return nil
		}
		print(log_error, conn.RemoteAddr, err)
		noticeError(conn, err.Error())
		return err
	}

	print(log_debug, "sent")

	return nil
}

func noticeAppendBroadcast(block *block_T) {
	data := struct {
		Hash string `json:"hash"`
		Index uint32 `json:"index"`
		StartTime int64 `json:"startTime"`
		Bits string `json:"bits"`
	} {}

	data.Hash = hex.EncodeToString(block.head.hash[:])
	data.Index = binary.LittleEndian.Uint32(block.head.hash[32:])
	data.StartTime = int64(binary.LittleEndian.Uint64(block.head.timestamp[:]))
	data.Bits = hex.EncodeToString(block.head.bits[:])

	noticeBroadcast(WS_NOTICE_APPEND, data)
}

func noticeTaskBroadcast(emit *taskEmit_T) {
	noticeBroadcast(WS_NOTICE_TASK, emit)
}

func noticeErrorBroadcast(err error) {
	for conn, _ := range noticeConns {
		noticeError(conn, err.Error())
	}
}

func noticeBroadcast(event int, data interface{}) {
	noticeData := noticeData_T { event, data }

	for conn, _ := range noticeConns {
		err := conn.WriteJSON(noticeData)
		if err != nil {
			if isWebsocketCloseError(err) {
				conn.Close()
				delete(noticeConns, conn)
				print(log_warning, conn.RemoteAddr(), err, "close normal closure")
				continue
			}
			print(log_error, conn.RemoteAddr, err)
			noticeError(conn, err.Error())
		}
	}

	print(log_debug, "batch sent")
}

func isWebsocketCloseError(err error) bool {
	if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) ||
	websocket.IsCloseError(err, websocket.CloseGoingAway) ||
	websocket.IsCloseError(err, websocket.CloseProtocolError) ||
	websocket.IsCloseError(err, websocket.CloseUnsupportedData) ||
	websocket.IsCloseError(err, websocket.CloseNoStatusReceived) ||
	websocket.IsCloseError(err, websocket.CloseAbnormalClosure) ||
	websocket.IsCloseError(err, websocket.CloseInvalidFramePayloadData) ||
	websocket.IsCloseError(err, websocket.ClosePolicyViolation) ||
	websocket.IsCloseError(err, websocket.CloseMessageTooBig) ||
	websocket.IsCloseError(err, websocket.CloseMandatoryExtension) ||
	websocket.IsCloseError(err, websocket.CloseInternalServerErr) ||
	websocket.IsCloseError(err, websocket.CloseServiceRestart) ||
	websocket.IsCloseError(err, websocket.CloseTryAgainLater) ||
	websocket.IsCloseError(err, websocket.CloseTLSHandshake) {
		return true
	}

	return false
}
