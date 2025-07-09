package main

import (
	"gopkg.in/ini.v1"
)

const __CONF__ = "config"

var conf *config

type config struct {
	remoteHosts []string // the remote p2p host addresses
}

func (c *config) read() error {
	cfg, err := ini.Load(__CONF__)
	if err != nil {
		return err
	}

	c.remoteHosts = cfg.Section(ini.DEFAULT_SECTION).Key("remote_hosts").Strings(",")
	return nil
}
