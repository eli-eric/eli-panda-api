package middlewares

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"panda/apigateway/helpers"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type IUserStatusValidator interface {
	ValidateUserEnabled(userUID string) (isEnabled bool, err error)
	InvalidateUser(userUID string)
	GetCacheEntries() []UserStatusCacheEntry
}

type UserStatusCacheEntry struct {
	UserUID   string    `json:"userUID"`
	IsEnabled bool      `json:"isEnabled"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type userStatusCacheEntry struct {
	isEnabled bool
	expiresAt time.Time
}

type UserStatusValidator struct {
	neo4jDriver *neo4j.Driver
	ttl         time.Duration
	cache       map[string]userStatusCacheEntry
	cacheMux    sync.RWMutex
}

func NewUserStatusValidator(driver *neo4j.Driver, ttlSeconds int) IUserStatusValidator {
	ttl := time.Duration(ttlSeconds) * time.Second
	if ttl <= 0 {
		ttl = time.Minute
	}

	return &UserStatusValidator{
		neo4jDriver: driver,
		ttl:         ttl,
		cache:       map[string]userStatusCacheEntry{},
	}
}

func (v *UserStatusValidator) ValidateUserEnabled(userUID string) (isEnabled bool, err error) {
	if userUID == "" {
		return false, fmt.Errorf("empty user uid")
	}

	now := time.Now()

	v.cacheMux.RLock()
	cacheItem, ok := v.cache[userUID]
	v.cacheMux.RUnlock()

	if ok && cacheItem.expiresAt.After(now) {
		return cacheItem.isEnabled, nil
	}

	session, err := helpers.NewNeo4jSession(*v.neo4jDriver)
	if err != nil {
		return false, err
	}

	dbEnabled, err := helpers.GetNeo4jSingleRecordSingleValue[bool](session, UserEnabledByUIDQuery(userUID))
	if err != nil {
		return false, err
	}

	v.cacheMux.Lock()
	v.cache[userUID] = userStatusCacheEntry{
		isEnabled: dbEnabled,
		expiresAt: now.Add(v.ttl),
	}
	v.cacheMux.Unlock()

	return dbEnabled, nil
}

func (v *UserStatusValidator) InvalidateUser(userUID string) {
	v.cacheMux.Lock()
	delete(v.cache, userUID)
	v.cacheMux.Unlock()
}

func (v *UserStatusValidator) GetCacheEntries() []UserStatusCacheEntry {
	now := time.Now()

	v.cacheMux.RLock()
	defer v.cacheMux.RUnlock()

	result := make([]UserStatusCacheEntry, 0, len(v.cache))
	for userUID, cacheEntry := range v.cache {
		if cacheEntry.expiresAt.Before(now) {
			// Skip expired entries; do not modify the cache during a read operation.
			continue
		}

		result = append(result, UserStatusCacheEntry{
			UserUID:   userUID,
			IsEnabled: cacheEntry.isEnabled,
			ExpiresAt: cacheEntry.expiresAt,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].UserUID < result[j].UserUID
	})

	return result
}

func UserEnabledByUIDQuery(userUID string) (result helpers.DatabaseQuery) {
	result.Query = `
	WITH $userUID AS userUID
	OPTIONAL MATCH (u:User {uid: userUID})
	RETURN CASE
		WHEN u IS NULL THEN false
		ELSE coalesce(u.isEnabled, false)
	END AS result`

	result.ReturnAlias = "result"
	result.Parameters = map[string]interface{}{
		"userUID": userUID,
	}

	return result
}
