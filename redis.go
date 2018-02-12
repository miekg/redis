package redis

import (
	"github.com/mediocregopher/radix.v2/pool"
)

func connect(addr string, amount int) (*pool.Pool, error) {
	p, err := pool.New("tcp", addr, amount)
	return p, err
}
