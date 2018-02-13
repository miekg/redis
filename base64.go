package redis

import (
	"encoding/base64"

	"github.com/miekg/dns"
)

// ignore errors for now.

func ToString(m *dns.Msg) string {
	b, _ := m.Pack()
	return base64.RawStdEncoding.EncodeToString(b)
}

func FromString(s string, ttl int) *dns.Msg {
	m := new(dns.Msg)
	b, _ := base64.RawStdEncoding.DecodeString(s)
	m.Unpack(b)

	msgTTL(m, ttl)

	return m
}
