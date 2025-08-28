package main

import (
	"errors"
)

type minedError struct {
	code string
	message string
}

func (e minedError) Error() string {
	return e.code + " Mined block error:" + e.message
}

type makeBlockError struct {
	code string
	message string
}

func (e makeBlockError) Error() string {
	return e.code + " Make block error:" + e.message
}

var errWrongType =  errors.New("Wrong type.")

var errSynchronizing = errors.New("synchronizing.")
var errPrevHashNotMatch = errors.New("The hash and prev hash are not match.")
var errTransactionRootNotMatch = errors.New("The transactionRoot and it's data are not match.")
var errStateRootNotMatch = errors.New("The stateRoot and it's data are not match.")
var errNonceExpired = errors.New("Nonce expired.")

var errZeroBlock = errors.New("Zero block")
var errBlockIdNotMatch = errors.New("block hash and index are not match.")
var errBlockHashFormat = errors.New("invalid block hash format.")
var errBlockHashing = errors.New("block hash and it's data are not match")
var errBits = errors.New("The bits are not match")
var errTransactionRootHash = errors.New("The transaction root and its hash root are not match")
var errStateRootHash = errors.New("The state and its hash root are not match")
var errStateRoot = errors.New("The state roots are not match")
