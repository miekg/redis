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

func FromString(s string) *dns.Msg {
	m := new(dns.Msg)
	b, _ := base64.RawStdEncoding.DecodeString(s)
	m.Unpack(b)
	return m
}
