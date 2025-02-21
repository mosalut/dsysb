package main

import (
	"flag"
	"strconv"
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
	logFile bool
}

var seedAddrs = make(map[string]bool)
var cmdFlag *cmdFlag_T

func init() {
	cmdFlag = &cmdFlag_T{}
	readFlags(cmdFlag)
	flag.Parse()
	cmdFlag.networkID = 0
}

func main() {
	showLogo()
	log.Println(*cmdFlag)

	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	}

	if cmdFlag.logFile {
		err := openLogFile(strconv.Itoa(cmdFlag.port))
		if err != nil {
			log.Fatal(err)
		}
	}

	peer = q2p.NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID)
	q2p.Set_connection_num(cmdFlag.cn)
	peer.TimeSendLost = 5
	peer.Timeout = 16
	peer.Successed = transportSuccessed
	peer.Failed = transportFailed

	print(log_debug, "peer:", peer)
	err := peer.Run()
	if err != nil {
		log.Fatal(err)
	}
	print(log_debug, "conn:", peer.Conn)

	initDB()
	initState()
	runHttpServer(cmdFlag.httpPort)
}

func readFlags(cmdFlag *cmdFlag_T) {
	flag.StringVar(&cmdFlag.ip, "ip", "0.0.0.0", "The P2P host IP")
	flag.IntVar(&cmdFlag.port, "port", 10000, "The P2P host Port")
	flag.StringVar(&cmdFlag.remoteHost, "remote_host", "", "Remote host address")
	flag.IntVar(&cmdFlag.cn, "cn", 32, "The max p2p connections")
	flag.StringVar(&cmdFlag.httpPort, "http_port", "20000", "HTTP run on")
	flag.BoolVar(&cmdFlag.logFile, "log_file", true, "Write log to file")
}
