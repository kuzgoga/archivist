package ai

import (
	"errors"
	"github.com/dgraph-io/badger/v4"
	"log"
)

type ChatProviderWithCache struct {
	provider ChatProvider
	db       *badger.DB
}

func NewChatProviderWithCache(provider ChatProvider) *ChatProviderWithCache {
	opts := badger.DefaultOptions("llmcache")
	opts.Logger = nil

	db, err := badger.Open(opts)
	if err != nil {
		log.Printf("Failed to open Badger DB: %v, falling back to uncached mode", err)
		return &ChatProviderWithCache{provider: provider}
	}

	return &ChatProviderWithCache{
		provider: provider,
		db:       db,
	}
}

func (c *ChatProviderWithCache) Ask(request string) (ChatResponse, error) {
	var cachedValue string
	err := c.db.View(func(txn *badger.Txn) error {
		item, e := txn.Get([]byte(request))
		if e != nil {
			return e
		}
		valCopy, e := item.ValueCopy(nil)
		if e != nil {
			return e
		}
		cachedValue = string(valCopy)
		return nil
	})

	if err == nil {
		return ChatResponse{
			Answer:     cachedValue,
			Successful: true,
		}, nil
	} else if errors.Is(err, badger.ErrKeyNotFound) {
		response, errProvider := c.provider.Ask(request)
		if errProvider != nil {
			return response, errProvider
		}

		err = c.db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(request), []byte(response.Answer))
		})

		if err != nil {
			log.Printf("Failed to cache response: %v", err)
		}

		return response, nil
	} else {
		log.Printf("LLM cache read error: %v", err)
		return c.provider.Ask(request)
	}
}

func (c *ChatProviderWithCache) Close() {
	if c.db != nil {
		if err := c.db.Close(); err != nil {
			log.Printf("Error closing Badger DB: %v", err)
		}
	}
}
