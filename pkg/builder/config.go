package builder

import "archivist/pkg/export"

type Config struct {
	Sources  UserSources `yaml:"sources" validate:"required"`
	Ai       AiProvider  `yaml:"ai" validate:"required"`
	Exporter export.Exporter
}
