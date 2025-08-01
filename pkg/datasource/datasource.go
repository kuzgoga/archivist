package datasource

type SourceItem struct {
	Tag   string
	Name  string
	Topic *string
}

type Datasource interface {
	GetPersons() []SourceItem
	GetTerms() []SourceItem
	GetDates() []SourceItem
}
