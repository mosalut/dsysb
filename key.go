package main

import (
	"bufio"
	"os"
	"fmt"
)

const sign = "dsysb>"

func keyEvent() {
	for {
		fmt.Printf(sign)
		reader := bufio.NewReader(os.Stdin)
		command, _ := reader.ReadString('\n')

		commandProcess(command[:len(command) - 1])
	}
}

func commandProcess(command string) {
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
	case "":
	default:
		fmt.Println("command:" + command + " not found")
	}
}
