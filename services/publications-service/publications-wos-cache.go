package publicationsservice

import (
	"sync"
	"time"

	"panda/apigateway/services/publications-service/models"
)

// ELIPANDA-501: Clarivate's Starter tier caps lookups per day, so a DOI that
// is previewed repeatedly (user opens the dialog twice, two editors look at
// the same paper) must not walk out to the upstream each time. A successful
// exact-record match is therefore cached in-process for a fixed window.
//
// In-process (not Neo4j-backed) deliberately: the quota is per API key, one
// gateway replica serves this deployment, and a cold restart costs at most
// one upstream call per DOI — the boring option that satisfies the budget.

const (
	wosCacheTTL        = 24 * time.Hour
	wosCacheMaxEntries = 512
)

type wosCacheEntry struct {
	hit       models.WosHit
	expiresAt time.Time
}

type wosLookupCache struct {
	mu      sync.Mutex
	entries map[string]wosCacheEntry
	ttl     time.Duration
	maxSize int
	now     func() time.Time
}

func newWosLookupCache(ttl time.Duration, maxSize int) *wosLookupCache {
	return &wosLookupCache{
		entries: make(map[string]wosCacheEntry),
		ttl:     ttl,
		maxSize: maxSize,
		now:     time.Now,
	}
}

// get returns a cached record only while it is fresh; expired entries are
// dropped on read so they never pin memory.
func (cache *wosLookupCache) get(doi string) (models.WosHit, bool) {
	if cache == nil {
		return models.WosHit{}, false
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()

	entry, found := cache.entries[doi]
	if !found {
		return models.WosHit{}, false
	}
	if cache.now().After(entry.expiresAt) {
		delete(cache.entries, doi)
		return models.WosHit{}, false
	}
	return entry.hit, true
}

// store keeps the entry set bounded: expired entries are evicted first, then
// the soonest-expiring ones (the oldest lookups) until the cap holds again.
func (cache *wosLookupCache) store(doi string, hit models.WosHit) {
	if cache == nil {
		return
	}
	cache.mu.Lock()
	defer cache.mu.Unlock()

	now := cache.now()
	cache.entries[doi] = wosCacheEntry{hit: hit, expiresAt: now.Add(cache.ttl)}

	if len(cache.entries) <= cache.maxSize {
		return
	}
	for key, entry := range cache.entries {
		if now.After(entry.expiresAt) {
			delete(cache.entries, key)
		}
	}
	for len(cache.entries) > cache.maxSize {
		oldestKey := ""
		var oldest time.Time
		for key, entry := range cache.entries {
			if oldestKey == "" || entry.expiresAt.Before(oldest) {
				oldestKey, oldest = key, entry.expiresAt
			}
		}
		delete(cache.entries, oldestKey)
	}
}
