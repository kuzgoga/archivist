package datasource

import (
	"errors"
	"github.com/fumiama/go-docx"
	"log"
	"os"
	"strings"
)

const (
	errInvalidTopicDelimiter = "docx: invalid topic delimiter"
	errNoParagraphsFound     = "docx: no paragraphs found or empty datasource"
	errUnknownTable          = "docx: fragment is unknown table"
	errEmptyCellInTable      = "docx: empty cell in table"
)

type SourcePosition struct {
	Filename          string
	Tag               string
	KeyPhrase         string
	ItemsDelimiter    string
	TrimSpaces        bool
	RemoveTrailingDot bool
	IsTable           bool
	TopicsDelimiter   *string
}

type DocxReader struct {
	personsSources []SourcePosition
	termsSources   []SourcePosition
	datesSources   []SourcePosition
	persons        []SourceItem
	terms          []SourceItem
	dates          []SourceItem
}

type Option func(*DocxReader)

func NewDocxReader(options ...Option) (*DocxReader, error) {
	reader := DocxReader{
		personsSources: nil,
		termsSources:   nil,
		datesSources:   nil,
		persons:        make([]SourceItem, 0),
		terms:          make([]SourceItem, 0),
		dates:          make([]SourceItem, 0),
	}

	for _, opt := range options {
		opt(&reader)
	}

	sources := map[string][]SourcePosition{}
	sources["persons"] = reader.personsSources
	sources["terms"] = reader.termsSources
	sources["dates"] = reader.datesSources

	buckets := map[string][]SourceItem{}
	buckets["persons"] = reader.persons
	buckets["terms"] = reader.terms
	buckets["dates"] = reader.dates

	for name, source := range sources {
		items, e := processSources(name, source)
		if e != nil {
			return nil, e
		}
		buckets[name] = items
	}

	reader.persons, reader.terms, reader.dates = buckets["persons"], buckets["terms"], buckets["dates"]
	return &reader, nil
}

func WithPersons(sources ...SourcePosition) Option {
	return func(reader *DocxReader) {
		reader.personsSources = sources
	}
}

func WithTerms(sources ...SourcePosition) Option {
	return func(reader *DocxReader) {
		reader.termsSources = sources
	}
}

func WithDates(sources ...SourcePosition) Option {
	return func(reader *DocxReader) {
		reader.datesSources = sources
	}
}

func (r *DocxReader) GetPersons() []SourceItem {
	return r.persons
}

func (r *DocxReader) GetDates() []SourceItem {
	return r.dates
}

func (r *DocxReader) GetTerms() []SourceItem {
	return r.terms
}

func processSources(sourceName string, sources []SourcePosition) ([]SourceItem, error) {
	var res []SourceItem
	if len(sources) != 0 {
		log.Printf("Parsing %s sources...", sourceName)
		for _, src := range sources {
			doc, err := openDocument(src.Filename)
			if err != nil {
				return nil, err
			}
			items, err := readSourceFromDocx(doc, src)
			if err != nil {
				return nil, err
			}
			res = append(res, items...)
		}
	} else {
		log.Printf("Skipping %s sources...\n", sourceName)
	}
	return res, nil
}

func openDocument(filename string) (*docx.Docx, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	doc, err := docx.Parse(file, fileInfo.Size())
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func readItemsFromTable(cell *docx.WTableCell, pos SourcePosition) ([]string, error) {
	if !pos.IsTable {
		return nil, errors.New(errUnknownTable)
	}

	paragraphs := make([]string, len(cell.Paragraphs))
	for i, p := range cell.Paragraphs {
		paragraphs[i] = p.String()
	}
	return paragraphs, nil
}

func readSourceFromDocx(doc *docx.Docx, pos SourcePosition) ([]SourceItem, error) {
	var paragraphs []string
	if !pos.IsTable {
		p := findParagraphByPhrase(doc, pos.KeyPhrase)
		if p == nil {
			return nil, errors.New(errNoParagraphsFound)
		}
		paragraphs = []string{
			p.String(),
		}
	} else {
		table := findTableByPhrase(doc, pos.KeyPhrase)
		if table == nil {
			return nil, errors.New(errEmptyCellInTable)
		}
		var err error
		paragraphs, err = readItemsFromTable(table, pos)
		if err != nil {
			return nil, err
		}
	}

	if len(paragraphs) == 0 {
		return nil, errors.New(errNoParagraphsFound)
	}

	return parseItemsFromParagraphs(pos, paragraphs)
}

func parseItemsFromParagraphs(pos SourcePosition, paragraphs []string) ([]SourceItem, error) {
	var items []SourceItem
	for _, p := range paragraphs {
		paragraphContent := p

		if pos.TopicsDelimiter != nil {
			topicAndElements := strings.SplitN(paragraphContent, *pos.TopicsDelimiter, 2)
			if len(topicAndElements) != 2 {
				return nil, errors.New(errInvalidTopicDelimiter)
			}

			topic, elements := topicAndElements[0], strings.Split(topicAndElements[1], pos.ItemsDelimiter)
			if pos.TrimSpaces {
				topic = strings.TrimSpace(topic)
			}

			for _, element := range elements {
				itemContent := element
				if pos.TrimSpaces {
					itemContent = strings.TrimSpace(itemContent)
				}
				if pos.RemoveTrailingDot {
					itemContent = strings.TrimSuffix(itemContent, ".")
				}
				items = append(items, SourceItem{
					Tag:   pos.Tag,
					Name:  itemContent,
					Topic: &topic,
				})
			}
		} else {
			elements := strings.Split(paragraphContent, pos.ItemsDelimiter)
			for _, element := range elements {
				itemContent := element
				if pos.TrimSpaces {
					itemContent = strings.TrimSpace(itemContent)
				}
				if pos.RemoveTrailingDot {
					itemContent = strings.TrimSuffix(itemContent, ".")
				}
				items = append(items, SourceItem{
					Tag:  pos.Tag,
					Name: itemContent,
				})
			}
		}
	}
	return items, nil
}

func findParagraphByPhrase(doc *docx.Docx, phrase string) *docx.Paragraph {
	for _, it := range doc.Document.Body.Items {
		switch it.(type) {
		case *docx.Paragraph:
			p := it.(*docx.Paragraph)
			if strings.Contains(p.String(), phrase) {
				return p
			}
		}
	}
	return nil
}

func findTableByPhrase(doc *docx.Docx, phrase string) *docx.WTableCell {
	for _, it := range doc.Document.Body.Items {
		switch it.(type) {
		case *docx.Table:
			t := it.(*docx.Table)
			for _, row := range t.TableRows {
				for _, cell := range row.TableCells {
					for _, p := range cell.Paragraphs {
						if strings.Contains(p.String(), phrase) {
							return cell
						}
					}
				}
			}
		}
	}
	return nil
}
