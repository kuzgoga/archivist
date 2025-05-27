package ai

import (
	"errors"
	"log"

	"github.com/dgraph-io/badger/v4"
)

type ChatProviderWithCache struct {
	provider ChatProvider
}

func NewChatProviderWithCache(provider ChatProvider) *ChatProviderWithCache {
	return &ChatProviderWithCache{
		provider: provider,
	}
}

func (c *ChatProviderWithCache) Ask(request string) (ChatResponse, error) {
	badgerOptions := badger.DefaultOptions("llmcache")
	badgerOptions.Logger = nil
	db, err := badger.Open(badgerOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var cachedValue string
	err = db.View(func(txn *badger.Txn) error {
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
		err = db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(request), []byte(response.Answer))
		})
		if err != nil {
			log.Fatal(err)
		}
		return response, nil
	} else {
		log.Fatal(err)
	}
	return ChatResponse{}, err
}
