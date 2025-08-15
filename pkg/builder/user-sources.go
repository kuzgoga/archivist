package builder

import "archivist/pkg/datasource"

type UserSources struct {
	Dates   []datasource.SourcePosition `yaml:"dates"`
	Terms   []datasource.SourcePosition `yaml:"terms"`
	Persons []datasource.SourcePosition `yaml:"persons"`
}
