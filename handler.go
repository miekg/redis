package redis

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
	"golang.org/x/net/context"
)

// ServeDNS implements the plugin.Handler interface.
func (re *Redis) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	qname := state.Name()
	qtype := state.QType()
	zone := plugin.Zones(re.Zones).Matches(qname)
	if zone == "" {
		return plugin.NextOrFailure(re.Name(), re.Next, ctx, w, r)
	}

	do := state.Do()

	now := re.now().UTC()

	m := re.get(now, qname, qtype, do)

	if m == nil {
		crr := &ResponseWriter{ResponseWriter: w, Redis: re}
		return plugin.NextOrFailure(re.Name(), re.Next, ctx, crr, r)
	}

	m.Id = r.Id // Copy IDs so the client will accept this answer.
	state.SizeAndDo(m)
	m, _ = state.Scrub(m)
	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (c *Redis) Name() string { return "redis" }
