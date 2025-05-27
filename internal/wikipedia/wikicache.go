package wikipedia

import (
	"errors"
	"github.com/dgraph-io/badger/v4"
	"log"
)

type ParserWithCache struct {
	parser WikiParser
}

func NewParserWithCache(parser WikiParser) *ParserWithCache {
	return &ParserWithCache{
		parser: parser,
	}
}

func (p *ParserWithCache) GetSummary(definition string) (string, error) {
	badgerOptions := badger.DefaultOptions("wikicache")
	badgerOptions.Logger = nil
	db, err := badger.Open(badgerOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	var cachedValue string

	err = db.View(func(txn *badger.Txn) error {
		item, e := txn.Get([]byte(definition))
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

	if errors.Is(err, badger.ErrKeyNotFound) {
		summary, err := p.parser.GetSummary(definition)
		if err != nil {
			return summary, err
		}
		err = db.Update(func(txn *badger.Txn) error {
			return txn.Set([]byte(definition), []byte(summary))
		})
		if err != nil {
			log.Fatal(err)
		}
		return summary, nil
	} else if err != nil {
		log.Fatal(err)
	}

	return cachedValue, nil
}
