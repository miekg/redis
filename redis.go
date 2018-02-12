package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/mediocregopher/radix.v2/pool"

	"github.com/miekg/dns"
)

// Redis is plugin that looks up responses in a cache and caches replies.
// It has a success and a denial of existence cache.
type Redis struct {
	Next  plugin.Handler
	Zones []string

	pool *pool.Pool
	nttl time.Duration
	pttl time.Duration

	addr string
	idle int
	// Testing.
	now func() time.Time
}

func New() *Redis {
	return &Redis{
		Zones: []string{"."},
		addr:  "localhost:6379",
		idle:  10,
		pool:  &pool.Pool{},
		pttl:  maxTTL,
		nttl:  maxNTTL,
		now:   time.Now,
	}
}

func Add(p *pool.Pool, key int, m *dns.Msg, duration time.Duration) error {
	// SETEX key duration m
	conn, err := p.Get()
	if err != nil {
		return err
	}
	defer p.Put(conn)

	println("adding", key)
	resp := conn.Cmd("SETEX", strconv.Itoa(key), duration.Seconds(), ToString(m))

	return resp.Err
}

func Get(p *pool.Pool, key int) (*dns.Msg, error) {
	println("GET", strconv.Itoa(key))
	// GET key
	conn, err := p.Get()
	if err != nil {
		println("no conn", err.Error())
		return nil, err
	}
	defer p.Put(conn)

	resp := conn.Cmd("GET", strconv.Itoa(key))
	if resp.Err != nil {
		println("err", resp.Err.Error())
		return nil, resp.Err
	}
	println("GET HERE")

	s, _ := resp.Str()
	if s == "" {
		println("string is empty")
		return nil, errors.New("not found")
	}

	m := FromString(s)

	println("GET RETURN", s, m.String())
	return m, nil
}

func (r *Redis) get(now time.Time, qname string, qtype uint16, do bool) *dns.Msg {
	k := hash(qname, qtype, do)

	if m, err := Get(r.pool, k); err != nil {
		println("hit", k)
		cacheHits.Inc()
		return m
	}
	println("miss", k)

	cacheMisses.Inc()
	return nil
}

func (r *Redis) connect() {
	// Can we ignore err here, i.e. will we try to connect later on?
	p, _ := pool.New("tcp", r.addr, r.idle)
	r.pool = p
	return
}
