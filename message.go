package main

import (
	"net"
	"fmt"

	"github.com/mosalut/q2p"
)

func sendMessage(peer *q2p.Peer_T, addr, message string) error {
	rAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	print(log_debug, message, rAddr)
	hash, err := peer.Transport(rAddr, []byte(message))
	if err != nil {
		return err
	}

	fmt.Println(hash, "sent")

	return nil
}
