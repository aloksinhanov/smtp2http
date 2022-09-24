package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/eko/gocache/cache"
	"github.com/eko/gocache/store"
	gocache "github.com/patrickmn/go-cache"
)

const (
	findAPIKey string = "SELECT cak.`key`FROM mw_customer_api_key cak INNER JOIN mw_customer cu ON cak.customer_id = cu.customer_id WHERE cak.`key`=? AND cu.email=?"
)

var (
	apiKeyAuth *APIKeyAuth
	authOnce   sync.Once
)

type APIKeyAuth struct {
	authdb    *DB
	authCache *cache.LoadableCache
}

type customerApiKey struct {
	Key string `db:"key"`
}

func NewAPIKeyAuthenticator(db *DB) *APIKeyAuth {
	authOnce.Do(func() {
		cacheClient := gocache.New(5*time.Minute, 1*time.Hour)
		cacheStore := store.NewGoCache(cacheClient, &store.Options{Expiration: 5 * time.Minute})
		apiKeyAuth = &APIKeyAuth{db, cache.NewLoadable(func(cacheApiKey interface{}) (interface{}, error) {

			log.Printf("\nNot in cache. Going to DB: %s", cacheApiKey)

			keys := make([]customerApiKey, 0)
			username, apikey := splitCacheKey(fmt.Sprintf("%s", cacheApiKey))

			err := db.FindAny(context.Background(), &keys, findAPIKey, []interface{}{apikey, username})
			if len(keys) == 0 {
				//not found, so do not cache this result
				return keys, fmt.Errorf("no record found for username: %s, apikey: %s", username, apikey)
			}

			return keys, err

		}, cacheStore)}
	})
	return apiKeyAuth
}

func makeCacheKey(username, apikey string) string {
	return fmt.Sprintf("%s\x00%s", username, apikey)
}

func splitCacheKey(cacheKey string) (username string, apikey string) {

	parts := strings.Split(cacheKey, "\x00")

	if len(parts) != 2 {
		//this should have never happened
		log.Printf("\nunexpected cache key: %+v", parts)
		panic(fmt.Errorf("unexpected cache key"))
	}

	username = parts[0]
	apikey = parts[1]

	return
}

func (a *APIKeyAuth) Authenticate(ctx context.Context, username, apiKey string) (bool, error) {

	var (
		err error
		ok  = false
	)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panicked %v", r)
			ok = false
		}
	}()

	keys, err := a.authCache.Get(makeCacheKey(username, apiKey))
	if err != nil {
		log.Printf("\nAuthenticate | FindAny: error: %v", err)
		return false, errors.New("failed to authenticate")
	}

	//if this panics, we have a rcover
	kk := keys.([]customerApiKey)

	if len(kk) == 0 {
		log.Println("Authenticate | FindAny: no keys found")
		return false, errors.New("access denied")
	}

	return ok, err
}
