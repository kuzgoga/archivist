package builder

type Config struct {
	Sources  UserSources `yaml:"sources" validate:"required"`
	Ai       AiProvider  `yaml:"ai" validate:"required"`
	Exporter Exporter    `yaml:"exporter"`
}
