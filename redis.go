package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

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

// New returns an new initialized Redis.
func New() *Redis {
	return &Redis{
		Zones: []string{"."},
		addr:  "127.0.0.1:6379",
		idle:  10,
		pool:  &pool.Pool{},
		pttl:  maxTTL,
		nttl:  maxNTTL,
		now:   time.Now,
	}
}

// Add adds the message m under k in Redis.
func Add(p *pool.Pool, key int, m *dns.Msg, duration time.Duration) error {
	// SETEX key duration m
	conn, err := p.Get()
	if err != nil {
		return err
	}
	defer p.Put(conn)

	resp := conn.Cmd("SETEX", strconv.Itoa(key), int(duration.Seconds()), ToString(m))

	return resp.Err
}

// Get returns the message under key from Redis.
func Get(p *pool.Pool, key int) (*dns.Msg, error) {
	conn, err := p.Get()
	if err != nil {
		return nil, err
	}
	defer p.Put(conn)

	resp := conn.Cmd("GET", strconv.Itoa(key))
	if resp.Err != nil {
		return nil, resp.Err
	}

	ttl := 0 // Item just expired, slap 0 TTL on it.
	respTTL := conn.Cmd("TTL", strconv.Itoa(key))
	if respTTL.Err == nil {
		ttl, err = respTTL.Int()
		if err != nil {
			ttl = 0
		}
	}

	s, _ := resp.Str()
	if s == "" {
		return nil, errors.New("not found")
	}

	m := FromString(s, ttl)

	return m, nil
}

func (re *Redis) get(now time.Time, state request.Request, server string) *dns.Msg {
	k := hash(state.Name(), state.QType(), state.Do())

	m, err := Get(re.pool, k)
	if err != nil {
		log.Debugf("Failed to get response from Redis cache: %s", err)
		cacheMisses.WithLabelValues(server).Inc()
		return nil
	}
	log.Debugf("Returning response from Redis cache: %s for %s", m.Question[0].Name, state.Name())
	cacheHits.WithLabelValues(server).Inc()
	return m
}

func (re *Redis) connect() (err error) {
	re.pool, err = pool.New("tcp", re.addr, re.idle)
	return err
}
