package redis

import (
	"encoding/base64"

	"github.com/miekg/dns"
)

// ToString converts the message m to a bsae64 encoded string.
func ToString(m *dns.Msg) string {
	b, _ := m.Pack()
	return base64.RawStdEncoding.EncodeToString(b)
}

// FromString converts s back into a DNS message.
func FromString(s string, ttl int) *dns.Msg {
	m := new(dns.Msg)
	b, _ := base64.RawStdEncoding.DecodeString(s)
	m.Unpack(b)

	msgTTL(m, ttl)

	return m
}
