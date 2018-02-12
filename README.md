# redis

## Name

*redis* - enables a networked cache.

## Description

With *redis* responses can be cached for up to 3600s. Caching in Redis is mostly usefull in
a setup where multiple CoreDNS instances share a VIP. E.g. multiple CoreDNS pods in a Kubernetes
cluster.

If Redis is not reacheable this plugin will be a noop. The *cache* and *redis* plugin can be used
together, where *cache* is the L1 and *redis* is the L2 level cache.
If multiple CoreDNS instances get a cache miss for the same item, they will all be fetching the same
information from an upstream and updating the cache, i.e. there is no (extra) coordination between
those instances.

There are two caches, one for positive and another for negative responses.

## Syntax

~~~ txt
redis [TTL] [ZONES...]
~~~

* **TTL** max TTL in seconds. If not specified, the maximum TTL will be used, which is 3600 for
    noerror responses and 1800 for denial of existence ones.
    Setting a TTL of 300: `cache 300` would cache records up to 300 seconds.
* **ZONES** zones it should cache for. If empty, the zones from the configuration block are used.

Each element in the cache is cached according to its TTL (with **TTL** as the max). For the negative
cache, the SOA's MinTTL value is used. A TTL of zero is not allowed.

If you want more control:

~~~ txt
redis [TTL] [ZONES...] {
    endpoint ENDPOINT...
    success TTL
    denial TTL
}
~~~

* **TTL**  and **ZONES** as above.
* `endpoint` specifies which **ENDPOINT** to use for Redis, this default to `localhost:6379`.
* `success`, override the settings for caching successful responses. **CAPACITY** indicates the maximum
  number of packets we cache before we start evicting (*randomly*). **TTL** overrides the cache maximum TTL.
* `denial`, override the settings for caching denial of existence responses. **CAPACITY** indicates the maximum
  number of packets we cache before we start evicting (LRU). **TTL** overrides the cache maximum TTL.
  There is a third category (`error`) but those responses are never cached.

## Metrics

If monitoring is enabled (via the *prometheus* directive) then the following metrics are exported:

* `coredns_redis_hits_total{type}` - Counter of cache hits by cache type.
* `coredns_redis_misses_total{}` - Counter of cache misses.

Cache types are either "denial" or "success".

## Examples

Enable caching for all zones, cache locally and also cache for up to 40s in the cluster wide Redis.

~~~ corefile
. {
    cache 30
    redis 40
    whoami
}
~~~

Proxy to Google Public DNS and only cache responses for example.org (or below).

~~~ corefile
. {
    proxy . 8.8.8.8:53
    redis example.org
}
~~~

# See Also

See [for more information](https://redis.io) on Redis.
