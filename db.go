package main

import (
	"strconv"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

var chainDB *leveldb.DB

func initDB() {
	var err error

	chainDB, err = leveldb.OpenFile("chain_" + strconv.Itoa(cmdFlag.networkID) + ".db", nil)
	if err != nil {
		log.Fatal(err)
	}
}
