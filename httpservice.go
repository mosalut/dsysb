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
	http.HandleFunc("/state2", state2Handler)
	http.HandleFunc("/assets", assetsHandler)
	http.HandleFunc("/asset2", asset2Handler)
	http.HandleFunc("/accounts", accountsHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/asset", assetHandler)
	http.HandleFunc("/account", accountHandler)
	http.HandleFunc("/account2", account2Handler)
	http.HandleFunc("/task", taskHandler)
	http.HandleFunc("/task2", task2Handler)
	http.HandleFunc("/blockchain", blockchainHandler)
	http.HandleFunc("/blockchain2", blockchain2Handler)
	http.HandleFunc("/block", blockHandler)
	http.HandleFunc("/block2", block2Handler)
	http.HandleFunc("/sendrawtransaction", sendRawTransactionHandler)
	http.HandleFunc("/transaction", getTransactionHandler)
	http.HandleFunc("/transaction2", getTransaction2Handler)

	http.HandleFunc("/txpool", txPoolHandler)

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
