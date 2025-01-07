package main

import (
	"crypto/sha256"
	"fmt"
	"log"
)

type script_T struct {
	creator *account_T
	byteCode []byte
}

func (script *script_T)hash(byteCode []byte) {
}

type vm_T struct {
	version string
}

func (vm *vm_T)write() {
	// TODO
}

func (vm *vm_T)read(byteCode []byte) {
	if len(byteCode) <= 20 {
		print("byte code length must > 20")
		return
	}

	script := &script_T{}
	script.byteCode = byteCode
	script.creator = &account_T {*address_T(byteCode[:16]), 0}
}
