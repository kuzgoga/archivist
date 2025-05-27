package wikipedia

import (
	"errors"
	"fmt"
	"github.com/trietmn/go-wiki"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type WikiParser interface {
	GetSummary(definition string) (string, error)
}

type Parser struct {
	summarySentences int
}

const UserAgent string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"
const language string = "ru"

func NewParser(summarySentencesAmount int) *Parser {
	gowiki.SetLanguage(language)
	gowiki.SetUserAgent(UserAgent)
	if summarySentencesAmount > 10 {
		panic("wikipedia: summarySentences amount > 10")
	}
	return &Parser{summarySentences: summarySentencesAmount}
}

func (p *Parser) GetSummary(definition string) (string, error) {
	results, _, err := gowiki.Search(definition, 10, true)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", errors.New("no results found")
	}
	for _, result := range results {
		fmt.Println(result)
	}
	summary, err := gowiki.Summary(results[0], p.summarySentences, 0, true, true)
	if err != nil {
		return "", err
	}
	return removeAccents(summary), nil
}

func removeAccents(s string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.Predicate(func(r rune) bool {
			return r == '\u0301' || r == '\u0300'
		})),
		norm.NFC,
	)

	result, _, err := transform.String(t, s)
	if err != nil {
		return s
	}

	return result
}
