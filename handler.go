package redis

import (
	"context"

	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/metrics"
	"github.com/coredns/coredns/request"

	"github.com/miekg/dns"
)

// ServeDNS implements the plugin.Handler interface.
func (re *Redis) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	zone := plugin.Zones(re.Zones).Matches(state.Name())
	if zone == "" {
		return plugin.NextOrFailure(re.Name(), re.Next, ctx, w, r)
	}

	server := metrics.WithServer(ctx)
	now := re.now().UTC()

	m := re.get(now, state, server)
	if m == nil {
		crr := &ResponseWriter{ResponseWriter: w, Redis: re, state: state, server: metrics.WithServer(ctx)}
		return plugin.NextOrFailure(re.Name(), re.Next, ctx, crr, r)
	}

	m.SetReply(r)
	w.WriteMsg(m)

	return dns.RcodeSuccess, nil
}

// Name implements the Handler interface.
func (re *Redis) Name() string { return "redisc" }
