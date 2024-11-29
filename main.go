package main

import (
	"flag"
	"log"

	"github.com/mosalut/q2p"
)

type cmdFlag_T struct {
	ip string
	port int
	remoteHost string
	networkID uint16
	cn int
	httpPort string
}

var seedAddrs = make(map[string]bool)
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

	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	peer = q2p.NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID, callback)
	q2p.Set_connection_num(cmdFlag.cn)
	err := openLogFile()
	if err != nil {
		log.Fatal(err)
	}

	print(0, "peer:", peer)
	err = peer.Run()
	if err != nil {
		log.Fatal(err)
	}
	print(0, "conn:", peer.Conn)

	runHttpServer(cmdFlag.httpPort)
}

func readFlags(cmdFlag *cmdFlag_T) {
	flag.StringVar(&cmdFlag.ip, "ip", "0.0.0.0", "UDP host IP")
	flag.IntVar(&cmdFlag.port, "port", 10000, "UDP host Port")
	flag.StringVar(&cmdFlag.remoteHost, "remote_host", "", "remote host address")
	flag.IntVar(&cmdFlag.cn, "cn", 32, "connection_num")
	flag.StringVar(&cmdFlag.httpPort, "http_port", "20000", "http run on")
}
