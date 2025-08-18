package builder

import (
	"archivist/pkg/ai"
	"archivist/pkg/ai/copilot"
	"archivist/pkg/ai/gigachat"
	"archivist/pkg/ai/openai"
	"archivist/pkg/datasource"
	"archivist/pkg/export"
	"archivist/pkg/pipeline"
	"io"
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type Application struct {
	DataSource datasource.Datasource
	AiProvider ai.ChatProvider
	Exporter   export.Exporter
}

func LoadConfig(filename string) *Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open config file: %v\n", err)
	}
	defer func() { _ = file.Close() }()

	var config Config
	dec := yaml.NewDecoder(
		file,
		yaml.Validator(validator.New()),
		yaml.Strict(),
	)
	err = dec.Decode(&config)
	if err != nil {
		log.Fatal(yaml.FormatError(err, true, true))
	}

	return &config
}

func CreateAiProvider(config *Config) ai.ChatProvider {
	var providersCount int
	var provider ai.ChatProvider
	var err error

	if config.Ai.CopilotSettings != nil {
		providersCount++
		provider, err = copilot.NewClient()
		if err != nil {
			log.Fatalf("failed to create Copilot client: %v\n", err)
		}
	}

	if config.Ai.GigachatSettings != nil {
		providersCount++
		provider, err = gigachat.NewGigaChat(
			config.Ai.GigachatSettings.ClientId,
			config.Ai.GigachatSettings.ClientSecret,
			config.Ai.GigachatSettings.Model,
		)
		if err != nil {
			log.Fatalf("failed to create GigaChat client: %v\n", err)
		}
	}

	if config.Ai.OpenAiSettings != nil {
		providersCount++

		provider = openai.NewClient(
			config.Ai.OpenAiSettings.ApiKey,
			config.Ai.OpenAiSettings.Model,
			config.Ai.OpenAiSettings.BaseUrl,
		)
	}

	if providersCount == 0 {
		log.Fatalln("no AI provider configured")
	}

	if providersCount > 1 {
		log.Fatalln("multiple AI providers configured, only one is allowed")
	}

	return provider
}

func CreateDataSource(config *Config) datasource.Datasource {
	var sources []datasource.Option
	sourcesCount := 0

	for _, source := range config.Sources.Persons {
		sources = append(sources, datasource.WithPersons(source))
		sourcesCount++
	}

	for _, source := range config.Sources.Terms {
		sources = append(sources, datasource.WithTerms(source))
		sourcesCount++
	}

	for _, source := range config.Sources.Dates {
		sources = append(sources, datasource.WithDates(source))
		sourcesCount++
	}

	if sourcesCount == 0 {
		log.Fatalln("no data sources configured")
	}

	ds, err := datasource.NewDocxReader(sources...)
	if err != nil {
		log.Fatalf("failed to create data source: %v\n", err)
	}
	return ds
}

func BuildApplication(configFile string) *Application {
	config := LoadConfig(configFile)

	fillDefaultValues(config)

	aiProvider := CreateAiProvider(config)
	dataSource := CreateDataSource(config)
	exporter := export.CreatePdfExporter(
		config.Exporter.PersonsParts,
		config.Exporter.TermsParts,
		config.Exporter.DatesParts,
		*config.Exporter.DeleteTypstFiles,
	)

	return &Application{
		DataSource: dataSource,
		AiProvider: aiProvider,
		Exporter:   exporter,
	}
}

func (app *Application) Close() {
	if closer, ok := app.AiProvider.(io.Closer); ok {
		_ = closer.Close()
	}

	if closer, ok := app.DataSource.(io.Closer); ok {
		_ = closer.Close()
	}

	if closer, ok := app.Exporter.(io.Closer); ok {
		_ = closer.Close()
	}
}

func (app *Application) Run() error {
	llmCached := ai.NewChatProviderWithCache(app.AiProvider)
	defer llmCached.Close()

	result := pipeline.ProcessDatasourceItems(app.DataSource, llmCached)

	return app.Exporter.Export(result)
}
