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
	networkID int
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

	conf = &config{}

	initTargetValues()
}

func main() {
	err := conf.read()
	if err != nil {
		log.Fatal(err)
	}

	showLogo()
	log.Println(*cmdFlag)
	log.Println(stdBlockNum)

	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	} else {
		for _, v := range conf.remoteHosts {
			seedAddrs[v] = false
		}

		conf.remoteHosts = nil
	}

	if cmdFlag.logFile {
		err := openLogFile(strconv.Itoa(cmdFlag.port))
		if err != nil {
			log.Fatal(err)
		}
	}

	peer = q2p.NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, uint16(cmdFlag.networkID))
	q2p.SetConnectionNum(cmdFlag.cn)
	peer.TimeSendLost = 5
	peer.Timeout = 16
	peer.LifeCycle = lifeCycle
	peer.Successed = transportSuccessed
	peer.Failed = transportFailed

	initDB()
	initIndex()

	err = peer.Run()
	if err != nil {
		log.Fatal(err)
	}

	runHttpServer(cmdFlag.httpPort)
}

func readFlags(cmdFlag *cmdFlag_T) {
	flag.StringVar(&cmdFlag.ip, "ip", "0.0.0.0", "The P2P host IP")
	flag.IntVar(&cmdFlag.port, "port", 10000, "The P2P host Port")
	flag.StringVar(&cmdFlag.remoteHost, "remote_host", "", "Remote host address")
	flag.IntVar(&cmdFlag.networkID, "network_id", 0, "The network_id: 0:mainnet 0x1~0x10:testnet 0x10:dev")
	flag.IntVar(&cmdFlag.cn, "connections", 32, "The max p2p connections")
	flag.StringVar(&cmdFlag.httpPort, "http_port", "20000", "HTTP run on")
	flag.BoolVar(&cmdFlag.logFile, "log_file", false, "Write log to file")
}
