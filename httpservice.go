package main

import (
	"net/http"
	"time"
	"log"
)

func runHttpServer(port string) {
	http.HandleFunc("/peer", peerHandler)
	http.HandleFunc("/socket", socketHandler)
	http.HandleFunc("/noticesocket", noticeSocketHandler)
	http.HandleFunc("/state", stateHandler)
	http.HandleFunc("/assets", assetsHandler)
	http.HandleFunc("/accounts", accountsHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/asset", assetHandler)
	http.HandleFunc("/account", accountHandler)
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/blockchain", blockchainHandler)
	http.HandleFunc("/block", blockHandler)
	http.HandleFunc("/sendrawtransaction", sendRawTransactionHandler)
	http.HandleFunc("/transaction", getTransactionHandler)

	http.HandleFunc("/txpool", txPoolHandler)
	http.HandleFunc("/txinpool", txInPoolHandler)

	http.HandleFunc("/message", sendMessageHandler)
	http.HandleFunc("/broadcast", broadcastHandler)

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
