package config

import (
	"go.uber.org/zap"
)

type GlobalConfig struct {
	Configs map[int]*UserConfig `json:"configs"`
	log     *zap.SugaredLogger
	raw     []byte
}

type UserConfig struct {
	Bypass   []string           `json:"Bypass" validate:"dive,ipv4|cidrv4|ipv6|cidrv6"`
	TProxies map[string]*TProxy `json:"tproxies" validate:"dive"`
	Rules    []Rule             `json:"rules" validate:"dive"`
}

type Rule struct {
	// Match is an regex expression
	// to match an cgroup path relative to the root of cgroupfs.
	Match string `json:"match" validate:"required"`

	// TProxy means that the traffic comes from this cgroup
	// should be redirected to a TPROXY server.s
	TProxy string `json:"tproxy" validate:"required_without_all=Drop Direct,excluded_with=Drop Direct"`
	// Drop means that the traffic comes from this cgroup will be dropped.
	Drop bool `json:"drop" validate:"required_without_all=TProxy Direct,excluded_with=TProxy Direct"`
	// Direct means that the traffic comes from this cgroup will not be touched.
	Direct bool `json:"direct" validate:"required_without_all=TProxy Drop,excluded_with=TProxy Drop"`
}

// TProxy describes a TPROXY server.
type TProxy struct {
	Name   string `json:"name"`
	NoUDP  bool   `json:"no-udp"`
	NoIPv6 bool   `json:"no-ipv6"`
	Port   uint16 `json:"port" validate:"required"`
	// Mark is the fire wall mark used to identify the TPROXY server
	// and trigger reroute operation of netfliter
	// from OUTPUT to PREROUTING internally.
	// It **NOT** means that this TPROXY server
	// must send traffic with the fwmark.
	// This mark is designed to be changeable for user
	// to make sure this mark is not conflict
	// with any fire wall mark in use.
	Mark FireWallMark `json:"mark" validate:"required"`
	// DNSHijack will hijack the dns request traffic
	// should redirect to this TPROXY server,
	// and send them to directory to a dns server described in DNSHijack.
	// This option is for fake-ip.
	DNSHijack *DNSHijack `json:"dns-hijack"`
}

type DNSHijack struct {
	IP   *string `json:"ip" validate:"ip4_addr"`
	Port uint16  `json:"port"`
	// If TCP is set to true,
	// tcp traffic will be hijacked, too,
	// when it's destination port is 53.
	TCP bool `json:"tcp"`
}

type FireWallMark uint32

type Opt func(c *GlobalConfig) (ret *GlobalConfig, err error)
