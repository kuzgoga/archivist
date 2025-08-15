package builder

type CopilotSettings struct {
	Model string `yaml:"model" validate:"required"`
}

type GigachatSettings struct {
	ClientId     string `yaml:"clientId" validate:"required"`
	ClientSecret string `yaml:"clientSecret" validate:"required"`
	Model        string `yaml:"model" validate:"required"`
}

type OpenAiSettings struct {
	Model   string  `yaml:"model" validate:"required"`
	ApiKey  string  `yaml:"apiKey" validate:"required"`
	BaseUrl *string `yaml:"baseUrl,omitempty"`
}

type AiProvider struct {
	CopilotSettings  *CopilotSettings  `yaml:"copilot,omitempty"`
	GigachatSettings *GigachatSettings `yaml:"gigachat,omitempty"`
	OpenAiSettings   *OpenAiSettings   `yaml:"openai,omitempty"`
}
