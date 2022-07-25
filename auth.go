package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

const (
	findAPIKey string = "SELECT `key` FROM mw_customer_api_key WHERE `key` = ?"
)

var (
	apiKeyAuth *APIKeyAuth
	authOnce   sync.Once
)

type APIKeyAuth struct {
	authdb *DB
}

type customerApiKey struct {
	Key string `db:"key"`
}

func NewAPIKeyAuthenticator(db *DB) *APIKeyAuth {
	authOnce.Do(func() {
		apiKeyAuth = &APIKeyAuth{db}
	})
	return apiKeyAuth
}

func (a *APIKeyAuth) Authenticate(ctx context.Context, apiKey string) (bool, error) {

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

	keys := make([]customerApiKey, 0)

	err = a.authdb.FindAny(ctx, &keys, findAPIKey, []interface{}{apiKey})
	if err != nil {
		log.Printf("\nAuthenticate | FindAny: error: %v", err)
		return false, errors.New("failed to check for authentication")
	}

	if len(keys) == 0 {
		log.Println("Authenticate | FindAny: no keys found")
		return false, errors.New("access denied")
	}

	return ok, err
}
