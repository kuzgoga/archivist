package main

import (
	"bismark/internal/ai"
	"bismark/internal/ai/copilot"
	"bismark/internal/datasource"
	"bismark/internal/export"
	"bismark/internal/pipeline"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

func main() {
	var topicDelimiter = " - "

	ds, err := datasource.NewDocxReader(
		datasource.WithPersons(
			datasource.SourcePosition{
				Filename:          "memo_history.docx",
				Tag:               "Всеобщая история. 10 класс. Персоналии",
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
				Filename:          "memo_history.docx",
				Tag:               "Всеобщая история. 10 класс. Хронология",
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
				Filename:          "memo_history.docx",
				Tag:               "Всеобщая история. 10 класс. Термины",
				KeyPhrase:         "I Балканская война",
				ItemsDelimiter:    ",",
				TrimSpaces:        true,
				RemoveTrailingDot: true,
				IsTable:           true,
				TopicsDelimiter:   &topicDelimiter,
			},
		),

		datasource.WithPersons(
			datasource.SourcePosition{
				Filename:          "memo_history.docx",
				Tag:               "История России. 10 класс. Персоналии",
				KeyPhrase:         "Николай Андреев",
				ItemsDelimiter:    ",",
				TopicsDelimiter:   &topicDelimiter,
				TrimSpaces:        true,
				RemoveTrailingDot: true,
				IsTable:           true,
			},
		),

		datasource.WithDates(
			datasource.SourcePosition{
				Filename:          "memo_history.docx",
				Tag:               "История России. 10 класс. Даты",
				KeyPhrase:         "17 (30) июля 1914 года",
				ItemsDelimiter:    ",",
				TrimSpaces:        true,
				RemoveTrailingDot: true,
				IsTable:           true,
				TopicsDelimiter:   &topicDelimiter,
			},
		),

		datasource.WithTerms(
			datasource.SourcePosition{
				Filename:          "memo_history.docx",
				Tag:               "История России. 10 класс. Термины",
				KeyPhrase:         "аграрный вопрос",
				ItemsDelimiter:    ",",
				TrimSpaces:        true,
				RemoveTrailingDot: true,
				IsTable:           true,
				TopicsDelimiter:   &topicDelimiter,
			},
		),
	)

	if err != nil {
		log.Fatalln(err)
	}

	copilotApi, err := copilot.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	llmCached := ai.NewChatProviderWithCache(copilotApi)
	defer llmCached.Close()

	res := pipeline.ProcessDatasourceItems(ds, llmCached)

	err = export.ToPdf(res)

	if err != nil {
		log.Fatalf("Export error occured: %s\n", err)
	}
}
