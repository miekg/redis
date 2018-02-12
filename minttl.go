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

	minTTL := maxTTL
	for _, r := range append(m.Answer, m.Ns...) {
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
