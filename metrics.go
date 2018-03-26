package redis

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	cacheHits = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redisc",
		Name:      "hits_total",
		Help:      "The count of cache hits.",
	})

	cacheMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redisc",
		Name:      "misses_total",
		Help:      "The count of cache misses.",
	})

	cacheDrops = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redisc",
		Name:      "drops_total",
		Help:      "The number responses that are not cached, because the reply is malformed.",
	})

	redisErr = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redisc",
		Name:      "set_errors_total",
		Help:      "The count of errors when adding entries to redis.",
	})
)
