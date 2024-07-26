package client

import (
	"net"
	"time"
)

type Config struct {
	LocalAddr   net.Addr
	DialTimeout time.Duration
}
