package redis

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"

	"github.com/coredns/caddy"
)

var log = clog.NewWithPlugin("redisc")

func init() { plugin.Register("redisc", setup) }

func setup(c *caddy.Controller) error {
	re, err := parse(c)
	if err != nil {
		return plugin.Error("redisc", err)
	}
	if err := re.connect(); err != nil {
		log.Warningf("Failed to connect to Redis at %s: %s", re.addr, err)
	} else {
		log.Infof("Connected to Redis at %s", re.addr)
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		re.Next = next
		return re
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
				if len(args) < 1 {
					return nil, c.ArgErr()
				}
				pttl, err := strconv.Atoi(args[0])
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
				if len(args) < 1 {
					return nil, c.ArgErr()
				}
				nttl, err := strconv.Atoi(args[0])
				if err != nil {
					return nil, err
				}
				// Reserve 0 (and smaller for future things)
				if nttl <= 0 {
					return nil, fmt.Errorf("cache TTL can not be zero or negative: %d", nttl)
				}
				re.nttl = time.Duration(nttl) * time.Second
			case "endpoint":
				args := c.RemainingArgs()
				if len(args) < 1 {
					return nil, c.ArgErr()
				}
				h, _, err := net.SplitHostPort(args[0])
				if err != nil && strings.Contains(err.Error(), "missing port in address") {
					if x := net.ParseIP(args[0]); x == nil {
						return nil, fmt.Errorf("failed to parse IP: %s", args[0])
					}

					re.addr = net.JoinHostPort(args[0], "6379")
					continue
				}
				if err != nil {
					return nil, err
				}
				// h should be a valid IP
				if x := net.ParseIP(h); x == nil {
					return nil, fmt.Errorf("failed to parse IP: %s", h)
				}
				re.addr = args[0]

			default:
				return nil, c.ArgErr()
			}
		}

		for i := range origins {
			origins[i] = plugin.Host(origins[i]).NormalizeExact()[0]
		}
		re.Zones = origins

		return re, nil
	}

	return nil, nil
}
