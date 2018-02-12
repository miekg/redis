package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/mediocregopher/radix.v2/pool"

	"github.com/miekg/dns"
)

func connect(addr string, amount int) (*pool.Pool, error) {
	p, err := pool.New("tcp", addr, amount)
	return p, err
}

func Add(p *pool.Pool, key int, m *dns.Msg, duration time.Duration) error {
	// SETEX key duration m
	conn, err := p.Get()
	if err != nil {
		return err
	}
	defer p.Put(conn)

	resp := conn.Cmd("SETEX", strconv.Itoa(key), duration.Seconds(), ToString(m))
	return resp.Err
}

func Get(p *pool.Pool, key int) (*dns.Msg, error) {
	// GET key
	conn, err := p.Get()
	if err != nil {
		return nil, err
	}
	defer p.Put(conn)

	resp := conn.Cmd("GET", strconv.Itoa(key))
	if resp.Err != nil {
		return nil, resp.Err
	}

	s, _ := resp.Str()
	if s == "" {
		return nil, errors.New("not found")
	}

	m := FromString(s)

	return m, nil
}

func (r *Redis) get(now time.Time, qname string, qtype uint16, do bool) *dns.Msg {
	k := hash(qname, qtype, do)

	if m, err := Get(r.ncache, k); err != nil {
		cacheHits.WithLabelValues(Denial).Inc()
		return m
	}

	if m, err := Get(r.pcache, k); err != nil {
		cacheHits.WithLabelValues(Success).Inc()
		return m
	}

	cacheMisses.Inc()

	return nil
}
