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

func NewDocxReader(filename string, options ...Option) (*DocxReader, error) {
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

	if len(reader.personsSources) != 0 {
		for _, src := range reader.personsSources {
			items, err := readSourceFromDocx(doc, src)
			reader.persons = append(reader.persons, items...)
			if err != nil {
				return nil, err
			}
		}
	} else {
		log.Println("Skipping persons sources...")
	}

	if len(reader.datesSources) != 0 {
		for _, src := range reader.datesSources {
			items, err := readSourceFromDocx(doc, src)
			reader.dates = append(reader.dates, items...)
			if err != nil {
				return nil, err
			}
		}
	} else {
		log.Println("Skipping dates source...")
	}

	if len(reader.termsSources) != 0 {
		for _, src := range reader.termsSources {
			items, err := readSourceFromDocx(doc, src)
			reader.terms = append(reader.terms, items...)
			if err != nil {
				return nil, err
			}
		}
	} else {
		log.Println("Skipping term source...")
	}

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
