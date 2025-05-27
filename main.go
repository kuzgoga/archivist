package main

import (
	"bismark/internal/ai"
	"bismark/internal/datasource"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
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

	for _, p := range ds.GetPersons() {
		fmt.Printf("%s, %s, %s\n", p.Tag, p.Name, *p.Topic)
	}

	for _, p := range ds.GetDates() {
		fmt.Printf("%s, %s, %s\n", p.Tag, p.Name, *p.Topic)
	}

	for _, p := range ds.GetTerms() {
		fmt.Printf("%s, %s, %s\n", p.Tag, p.Name, *p.Topic)
	}

	gigachat, err := ai.NewGigaChat(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("GIGACHAT_MODEL"))

	if err != nil {
		panic(err)
	}

	llmCached := ai.NewChatProviderWithCache(gigachat)

	response, err := llmCached.Ask(`Дай историческую характеристику главным событиям, происходившим "8 сентября 1943 года". Дай одноабзацный ответ в формате "<дата> - <описание исторического события>" без разметки`)
	if err != nil {
		panic(err)
	}
	if response.Successful {
		fmt.Println("Successful")
	}
	fmt.Println(response.Answer)
}
