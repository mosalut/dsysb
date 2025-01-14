package main

import (
	"net/http"
	"time"
	"log"
)

func runHttpServer(port string) {
	http.HandleFunc("/peer", peerHandler)
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/blockchain", blockchainHandler)
	http.HandleFunc("/block", blockHandler)
	http.HandleFunc("/sendrawtransaction", sendRawTransaction)
	http.HandleFunc("/txpool", txPool)

	server := http.Server {
		Addr: "0.0.0.0:" + port,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	print(log_info, "HTTP server is running on 0.0.0.0:" + port)

	log.Fatal(server.ListenAndServe())
}

func cors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,UPDATE,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, X-Extra-Header, Content-Type, Accept, Authorization, id, username, mobile, token")
}
