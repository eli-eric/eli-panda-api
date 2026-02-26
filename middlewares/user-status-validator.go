package middlewares

import (
	"fmt"
	"sync"
	"time"

	"panda/apigateway/helpers"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type IUserStatusValidator interface {
	ValidateUserEnabled(userUID string) (isEnabled bool, err error)
	InvalidateUser(userUID string)
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

func UserEnabledByUIDQuery(userUID string) (result helpers.DatabaseQuery) {
	result.Query = `
	MATCH(u:User{uid: $userUID})
	RETURN coalesce(u.isEnabled, false) AS result`

	result.ReturnAlias = "result"
	result.Parameters = map[string]interface{}{
		"userUID": userUID,
	}

	return result
}
