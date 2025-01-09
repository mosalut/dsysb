package main

import (
	"gopkg.in/ini.v1"
	"log"
)

const __CONF__ = "config"

var conf *config

type config struct {
	version string
	connectionNum int // the max p2p connections
	remoteHost string // the remote p2p host address
	ip string // The P2P host IP
	port int // The P2P host port
	httpPort string // HTTP run on
}

func (c *config) read() error {
	cfg, err := ini.Load(__CONF__)
	if err != nil {
		return err
	}

	c.version = cfg.Section("").Key("version").String()
	c.connectionNum, err = cfg.Section("").Key("connection_num").Int()
	if err != nil {
		log.Fatal(err)
	}
	c.remoteHost = cfg.Section("").Key("remote_host").String()
	c.ip = cfg.Section("").Key("ip").String()
	c.port, err = cfg.Section("").Key("port").Int()
	if err != nil {
		log.Fatal(err)
	}
	c.httpPort = cfg.Section("").Key("http_port").String()
	return nil
}
