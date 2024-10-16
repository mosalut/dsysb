package main

import (
	"bufio"
	"os"
	"strings"
	"fmt"
)

const sign = "dsysb>"
const commandHelp = `
	conns: List connection addresses
	exit: exit
	quit: exit
	newaddress: generate a new address
	help: this doc
`

func keyEvent() {
	for {
		fmt.Printf(sign)
		reader := bufio.NewReader(os.Stdin)
		commandLine, _ := reader.ReadString('\n')

		commandProcess(commandLine)
	}
}

func commandProcess(commandLine string) {
	commandWords := strings.Fields(strings.TrimSpace(commandLine))
	command := commandWords[0]
	switch command {
	case "conns":
		for k, _ := range seedAddrs {
			fmt.Println(k)
		}
	case "set_conn_num":
	case "exit":
		os.Exit(0)
	case "quit":
		os.Exit(0)
	case "help":
		fmt.Println(commandHelp)
	case "newaddress":
		wallet := newWallet()
		address := wallet.getAddress()
		fmt.Println(string(address))
	case "validateaddress":
		fmt.Println(validateAddress(commandWords[1]))
	default:
		fmt.Println(commandHelp)
	}
}
