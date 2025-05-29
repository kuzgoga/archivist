package main

import (
	"bismark/internal/ai"
	"bismark/internal/ai/copilot"
	"bismark/internal/datasource"
	"bismark/internal/pipeline"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var topicDelimiter = " - "

	ds, err := datasource.NewDocxReader(
		datasource.WithPersons(
			datasource.SourcePosition{
				Filename:          "memo_history.docx",
				Tag:               "10ClassWorldPersons",
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
				Tag:               "10ClassWorldDates",
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
				Tag:               "10ClassWorldTerms",
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
				Tag:               "10ClassRuPersons",
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
				Tag:               "10ClassRuDates",
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
				Tag:               "10ClassRuTerms",
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
		panic(err)
	}

	//gigaChat := openai.NewClient(os.Getenv("OPENAI_API_KEY"), shared.ChatModelGPT4oMini)
	gigaChat, err := copilot.NewClient()
	if err != nil {
		panic(err)
	}
	llmCached := ai.NewChatProviderWithCache(gigaChat)

	res := pipeline.ProcessDatasourceItems(ds, llmCached)

	for _, item := range res.Persons {
		fmt.Printf("%s: %s: %s\n", item.Tag, item.Topic, item.Summary)
	}
	for _, item := range res.Dates {
		fmt.Printf("%s: %s: %s\n", item.Tag, item.Topic, item.Summary)
	}
	for _, item := range res.Terms {
		fmt.Printf("%s: %s: %s\n", item.Tag, item.Topic, item.Summary)
	}
}
