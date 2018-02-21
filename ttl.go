package redis

import (
	"time"

	"github.com/coredns/coredns/plugin/pkg/response"

	"github.com/miekg/dns"
)

func minMsgTTL(m *dns.Msg, mt response.Type) time.Duration {
	if mt != response.NoError && mt != response.NameError && mt != response.NoData {
		return 0
	}

	// No data to examine, return a short ttl as a fail safe.
	if len(m.Answer)+len(m.Ns) == 0 {
		return failSafeTTL
	}

	minTTL := maxTTL
	for _, r := range append(append(m.Answer, m.Ns...), m.Extra...) {
		if r.Header().Rrtype == dns.TypeOPT {
			// OPT records use TTL field for extended rcode and flags
			continue
		}
		switch mt {
		case response.NameError, response.NoData:
			if r.Header().Rrtype == dns.TypeSOA {
				return time.Duration(r.(*dns.SOA).Minttl) * time.Second
			}
		case response.NoError, response.Delegation:
			if r.Header().Ttl < uint32(minTTL.Seconds()) {
				minTTL = time.Duration(r.Header().Ttl) * time.Second
			}
		}
	}
	return minTTL
}

func msgTTL(m *dns.Msg, ttl int) {
	for i := range m.Answer {
		m.Answer[i].Header().Ttl = uint32(ttl)
	}
	for i := range m.Ns {
		m.Ns[i].Header().Ttl = uint32(ttl)
	}
	for i := range m.Extra {
		if m.Extra[i].Header().Rrtype != dns.TypeOPT {
			m.Extra[i].Header().Ttl = uint32(ttl)
		}
	}
}
