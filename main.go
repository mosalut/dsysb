package main

import (
	"syscall"
	"os"
	"os/signal"
	"flag"
	"net"
	"log"

	"github.com/mosalut/q2p"
)

type cmdFlag_T struct {
	ip string
	port int
	remoteHost string
	networkID uint16
}

var cmdFlag *cmdFlag_T

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cmdFlag = &cmdFlag_T{}
	readFlags(cmdFlag)
	flag.Parse()
	cmdFlag.networkID = 0
}

func main() {
	log.Println(*cmdFlag)

	seedAddrs := make(map[*net.UDPAddr]bool)
	if cmdFlag.remoteHost != "" {
		remoteAddr, err := net.ResolveUDPAddr("udp", cmdFlag.remoteHost)
		if err != nil {
			log.Fatal(err)
		}

		seedAddrs[remoteAddr] = false
	}

	peer := q2p.NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID)

	log.Println("peer:", peer)
	err := peer.Run()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("conn:", peer.Conn)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Println("Received signal, shutting down...")
}

func readFlags(cmdFlag *cmdFlag_T) {
	flag.StringVar(&cmdFlag.ip, "ip", "0.0.0.0", "UDP host IP")
	flag.IntVar(&cmdFlag.port, "port", 10000, "UDP host Port")
	flag.StringVar(&cmdFlag.remoteHost, "remote_host", "", "remote host address")
}
