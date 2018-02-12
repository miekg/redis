package redis

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"

	"github.com/mholt/caddy"
)

func init() {
	caddy.RegisterPlugin("redis", caddy.Plugin{
		ServerType: "dns",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	re, err := parse(c)
	if err != nil {
		return plugin.Error("redis", err)
	}
	re.connect()

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		re.Next = next
		return re
	})

	c.OnStartup(func() error {
		once.Do(func() {
			m := dnsserver.GetConfig(c).Handler("prometheus")
			if m == nil {
				return
			}
			if x, ok := m.(*metrics.Metrics); ok {
				x.MustRegister(cacheHits)
				x.MustRegister(cacheMisses)
			}
		})
		return nil
	})

	return nil
}

func parse(c *caddy.Controller) (*Redis, error) {
	re := New()

	for c.Next() {
		// cache [ttl] [zones..]
		origins := make([]string, len(c.ServerBlockKeys))
		copy(origins, c.ServerBlockKeys)
		args := c.RemainingArgs()

		if len(args) > 0 {
			// first args may be just a number, then it is the ttl, if not it is a zone
			ttl, err := strconv.Atoi(args[0])
			if err == nil {
				// Reserve 0 (and smaller for future things)
				if ttl <= 0 {
					return nil, fmt.Errorf("cache TTL can not be zero or negative: %d", ttl)
				}
				re.pttl = time.Duration(ttl) * time.Second
				re.nttl = time.Duration(ttl) * time.Second
				args = args[1:]
			}
			if len(args) > 0 {
				copy(origins, args)
			}
		}

		// Refinements? In an extra block.
		for c.NextBlock() {
			switch c.Val() {
			case Success:
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				pttl, err := strconv.Atoi(args[1])
				if err != nil {
					return nil, err
				}
				// Reserve 0 (and smaller for future things)
				if pttl <= 0 {
					return nil, fmt.Errorf("cache TTL can not be zero or negative: %d", pttl)
				}
				re.pttl = time.Duration(pttl) * time.Second
			case Denial:
				args := c.RemainingArgs()
				if len(args) != 1 {
					return nil, c.ArgErr()
				}
				nttl, err := strconv.Atoi(args[1])
				if err != nil {
					return nil, err
				}
				// Reserve 0 (and smaller for future things)
				if nttl <= 0 {
					return nil, fmt.Errorf("cache TTL can not be zero or negative: %d", nttl)
				}
				re.nttl = time.Duration(nttl) * time.Second
			case "endpoint":
				println("not implemented")

			default:
				return nil, c.ArgErr()
			}
		}

		for i := range origins {
			origins[i] = plugin.Host(origins[i]).Normalize()
		}
		re.Zones = origins

		return re, nil
	}

	return nil, nil
}

var once sync.Once
