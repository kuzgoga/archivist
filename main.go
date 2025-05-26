package main

import (
	"bismark/internal/datasource"
	"fmt"
)

// Not a table error
// Selectors

func main() {
	var topicDelimiter = " - "

	ds, err := datasource.NewDocxReader("memo_history.docx",
		datasource.WithPersons(
			datasource.SourcePosition{
				Tag:               "PersonsRu",
				KeyPhrase:         "Вильгельм II",
				ItemsDelimiter:    ",",
				TopicsDelimiter:   &topicDelimiter,
				TrimSpaces:        true,
				RemoveTrailingDot: true,
				IsTable:           true,
			},
		),

		datasource.WithDates(
			datasource.SourcePosition{
				Tag:               "Dates",
				KeyPhrase:         "28 июня 1914",
				ItemsDelimiter:    ",",
				TrimSpaces:        true,
				RemoveTrailingDot: true,
				IsTable:           true,
				TopicsDelimiter:   &topicDelimiter,
			},
		),
		
		datasource.WithTerms(
			datasource.SourcePosition{
				Tag:               "Terms",
				KeyPhrase:         "I Балканская война",
				ItemsDelimiter:    ",",
				TrimSpaces:        true,
				RemoveTrailingDot: true,
				IsTable:           true,
				TopicsDelimiter:   &topicDelimiter,
			},
		),
	)

	if err != nil {
		panic(err)
	}

	for _, p := range ds.GetPersons() {
		fmt.Printf("%s, %s, %s\n", p.Tag, p.Name, *p.Topic)
	}

	for _, p := range ds.GetDates() {
		fmt.Printf("%s, %s, %s\n", p.Tag, p.Name, *p.Topic)
	}

	for _, p := range ds.GetTerms() {
		fmt.Printf("%s, %s, %s\n", p.Tag, p.Name, *p.Topic)
	}

}
