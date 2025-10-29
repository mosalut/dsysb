package main

import (
	"flag"
	"strconv"
	"encoding/hex"
	"fmt"
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

	displayLogo()

	if cmdFlag.remoteHost != "" {
		seedAddrs[cmdFlag.remoteHost] = false
	} else {
		for _, v := range conf.remoteHosts {
			seedAddrs[v] = false
		}

		fmt.Println("2222:")
		for k, v := range seedAddrs {
			fmt.Println(k, v)
		}
		conf.remoteHosts = nil
		fmt.Println("3333:")
		for k, v := range seedAddrs {
			fmt.Println(k, v)
		}
	}

	displayInitInfo()

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

func displayInitInfo() {
	fmt.Println(*cmdFlag)
	fmt.Println("-------------------------------------------")
	fmt.Printf("Network ID:%d\n", cmdFlag.networkID)
	fmt.Println("P2P host on:", cmdFlag.ip + ":" + fmt.Sprintf("%d", cmdFlag.port))
	fmt.Println("The max of p2p connect:", cmdFlag.cn)
	fmt.Println("The http port:", cmdFlag.httpPort)
	fmt.Println("Block period batch:", stdBlockNum)
	fmt.Println("Remote hosts:")
	for k, v := range seedAddrs {
		fmt.Println("\t", k, v)
	}
	fmt.Println("The block period batch:", stdBlockNum)
	fmt.Println("The block period batch SECs:", stdBlockBatchSeconds)
	fmt.Println("The current period batch difficult:", hex.EncodeToString(difficult_1_target[:]))
	fmt.Println("-------------------------------------------")
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
