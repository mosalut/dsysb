package main

import (
	"net"
//	"log"

	"github.com/mosalut/q2p"
)

func sendMessage(peer *q2p.Peer_T, addr, message string) error {
	rAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	print(0, message, rAddr)
	err = peer.Transport(rAddr, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
