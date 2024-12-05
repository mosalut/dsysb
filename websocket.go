// server
package main

import (
	"net/http"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

const (
	TILL_BLOCK = iota
	NEW_BLOCK
)

type socketData_T struct {
	Event int `json:"event"`
	Body interface{} `json:"body"`
}

var upgrader = websocket.Upgrader{}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade our raw HTTP connection to a websocket based one
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		print(2, "Error during connection upgradation:", err)
		return
	}
//	defer conn.Close()

	if len(blockchain) != 0 {
		data := socketData_T {TILL_BLOCK, blockchain[len(blockchain) - 1].Head}
		err := conn.WriteJSON(data)
		if err != nil {
			print(2, err)
			return
		}
	}

	// The event loop
	for {
		data := socketData_T{}
		err := conn.ReadJSON(&data)
		if err != nil {
			print(2, err)
			continue
		}

		print(1, data)

		switch data.Event {
		case NEW_BLOCK:
			log.Println("new block")
			bm := data.Body.(map[string]interface{})
			for _, v := range bm {
				fmt.Println(v)
			}

			head := &blockHead_T {
				int(bm["head"].(map[string]interface{})["index"].(float64)),
				int64(bm["head"].(map[string]interface{})["timestamp"].(float64)),
				bm["head"].(map[string]interface{})["hash"].(string),
				bm["head"].(map[string]interface{})["prev_hash"].(string),
				int(bm["head"].(map[string]interface{})["nonce"].(float64)),
			}
			body := &blockBody_T {}
			block := &block_T {
				head,
				body,
			}


			blockchain = append(blockchain, block)
		}
	}
}
