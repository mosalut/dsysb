package main

import (
	"syscall"
	"os"
	"os/signal"
	"flag"
	"time"
	"log"

	"github.com/mosalut/q2p"
)

type cmdFlag_T struct {
	ip string
	port int
	remoteHost string
	networkID uint16
	cn int
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

	go keyEvent()

	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	peer := q2p.NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID)
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

	go func () {
		for {
			print(0, "xxxxxxxxxxxxxxxxxxxxxx")
			time.Sleep(time.Second)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	print(0, "Received signal, shutting down...")
}

func readFlags(cmdFlag *cmdFlag_T) {
	flag.StringVar(&cmdFlag.ip, "ip", "0.0.0.0", "UDP host IP")
	flag.IntVar(&cmdFlag.port, "port", 10000, "UDP host Port")
	flag.StringVar(&cmdFlag.remoteHost, "remote_host", "", "remote host address")
	flag.IntVar(&cmdFlag.cn, "cn", 32, "connection_num")
}
