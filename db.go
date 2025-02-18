package main

import (
	"strconv"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

var chainDB *leveldb.DB

func initDB() {
	port := strconv.Itoa(cmdFlag.port)

	var err error

	chainDB, err = leveldb.OpenFile("chain_" + port + ".db", nil)
	if err != nil {
		log.Fatal(err)
	}
}
