package builder

type Exporter struct {
	PersonsParts     int   `yaml:"personsParts"`
	TermsParts       int   `yaml:"termsParts"`
	DatesParts       int   `yaml:"datesParts"`
	DeleteTypstFiles *bool `yaml:"deleteTypstFiles"`
}
