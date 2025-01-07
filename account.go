package main

type address_T []byte

type account_T struct {
	address *address_T
	nonce int
}
