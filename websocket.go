// server
package main

import (
	"encoding/binary"
	"encoding/json"
	"net/http"
	"log"

	"github.com/gorilla/websocket"
)

const (
	TILL_BLOCK = iota
	NEW_BLOCK
)

type socketData_T struct {
	Event int `json:"event"`
	Body []byte `json:"body"`
}

var upgrader = websocket.Upgrader{}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		print(log_error, "Error during connection upgradation:", err)
		return
	}
//	defer conn.Close()

/*
	if len(blockchain) != 0 {
		data := socketData_T {TILL_BLOCK, blockchain[len(blockchain) - 1].Head}
		err := conn.WriteJSON(data)
		if err != nil {
			print(log_error, err)
			return
		}
	}
	*/

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

		print(log_info, data)

		switch data.Event {
		case NEW_BLOCK:
			log.Println("new block")

			block := &block_T{}
			err = json.Unmarshal(data.Body, block)
			if err != nil {
				print(log_error, err)
				break
			}
			log.Println(block.Head.PrevHash)

			buffer := make([]byte, 4, 4)
			binary.LittleEndian.PutUint32(buffer, block.Head.Index)

			err = toolDB.Put([]byte("till"), buffer, nil)
			if err != nil {
				print(log_error, err)
				break
			}

			err = chainDB.Put(buffer, data.Body, nil)
			if err != nil {
				print(log_error, err)
				break
			}
		}
	}
}
