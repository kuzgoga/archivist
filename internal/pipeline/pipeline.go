package pipeline

import (
	"bismark/internal/ai"
	"bismark/internal/datasource"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"log"
)

type CompleteItem struct {
	Tag     string
	Topic   string
	Name    string
	Summary string
}

type Result struct {
	Persons []CompleteItem
	Terms   []CompleteItem
	Dates   []CompleteItem
}

func ProcessDatasourceItems(datasource datasource.Datasource, llm ai.ChatProvider) Result {
	result := Result{
		Persons: nil,
		Terms:   nil,
		Dates:   nil,
	}
	result.Persons = processItemsList(datasource.GetPersons(), llm, ai.PersonPrompt, "persons")
	result.Terms = processItemsList(datasource.GetTerms(), llm, ai.TermPrompt, "terms")
	result.Dates = processItemsList(datasource.GetDates(), llm, ai.DatePrompt, "dates")
	return result
}

func processItemsList(source []datasource.SourceItem, llm ai.ChatProvider, prompt string, name string) []CompleteItem {
	var result []CompleteItem
	bar := progressbar.NewOptions(
		len(source),
		progressbar.OptionSetDescription(fmt.Sprintf("Writing %s...", name)),
	)

	for _, el := range source {
		p := fmt.Sprintf(prompt, el.Name)
		response, err := llm.Ask(p)
		if err != nil {
			log.Printf("LLM error for prompt `%s`: %s", p, err.Error())
			_ = bar.Add(1)
			continue
		}
		fmt.Println(response.Answer) // TODO: remove

		var topic string
		if el.Topic != nil {
			topic = *el.Topic
		} else {
			topic = ""
		}

		item := CompleteItem{
			Tag:     el.Tag,
			Topic:   topic,
			Name:    el.Name,
			Summary: response.Answer,
		}
		result = append(result, item)
		_ = bar.Add(1)
	}
	return result
}
