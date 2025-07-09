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
//	cmdFlag.networkID = 0x0 // 0:mainnet
//	cmdFlag.networkID = 0x1 // 0x1~0x10:testnet
	cmdFlag.networkID = 0x10 // 0x10:dev

	if cmdFlag.networkID == 0x10 {
		// dev
		stdBlockNum = 100 // for test faster
		stdBlockBatchSeconds = 60000 // 600 * 100 for dev faster
		difficult_1_target = [4]byte{ 0x1f, 0x00, 0xff, 0xff }

	} else {
		// others
		stdBlockNum = 1024
		stdBlockBatchSeconds = 614400 // 600 * 1024
		difficult_1_target = [4]byte{ 0x1d, 0, 0xff, 0xff }
	}

	conf = &config{}
}

func main() {
	err := conf.read()
	if err != nil {
		return
	}

	showLogo()
	log.Println(*cmdFlag)

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

	peer = q2p.NewPeer(cmdFlag.ip, cmdFlag.port, seedAddrs, cmdFlag.networkID)
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
	flag.IntVar(&cmdFlag.cn, "cn", 32, "The max p2p connections")
	flag.StringVar(&cmdFlag.httpPort, "http_port", "20000", "HTTP run on")
	flag.BoolVar(&cmdFlag.logFile, "log_file", false, "Write log to file")
}
