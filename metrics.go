package redis

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redis",
		Name:      "hits_total",
		Help:      "The count of cache hits.",
	})

	cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redis",
		Name:      "misses_total",
		Help:      "The count of cache misses.",
	})

	redisErr = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redis",
		Name:      "set_errors_total",
		Help:      "The count of errors when adding entries to redis.",
	})
)
