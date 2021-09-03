package redis

import (
	"github.com/coredns/coredns/plugin"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	cacheHits = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redisc",
		Name:      "hits_total",
		Help:      "The count of cache hits.",
	}, []string{"server"})

	cacheMisses = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redisc",
		Name:      "misses_total",
		Help:      "The count of cache misses.",
	}, []string{"server"})

	cacheDrops = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redisc",
		Name:      "drops_total",
		Help:      "The number responses that are not cached, because the reply is malformed.",
	}, []string{"server"})

	redisErr = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "redisc",
		Name:      "set_errors_total",
		Help:      "The count of errors when adding entries to redis.",
	}, []string{"server"})
)
