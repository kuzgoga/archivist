package builder

func fillDefaultValues(config *Config) {
	for i, source := range config.Sources.Persons {
		if source.TrimSpaces == nil {
			config.Sources.Persons[i].TrimSpaces = new(bool)
			*config.Sources.Persons[i].TrimSpaces = true
		}
		if source.RemoveTrailingDot == nil {
			config.Sources.Persons[i].RemoveTrailingDot = new(bool)
			*config.Sources.Persons[i].RemoveTrailingDot = true
		}
		if source.ItemsDelimiter == "" {
			config.Sources.Persons[i].ItemsDelimiter = ","
		}
	}

	for i, source := range config.Sources.Terms {
		if source.TrimSpaces == nil {
			config.Sources.Terms[i].TrimSpaces = new(bool)
			*config.Sources.Terms[i].TrimSpaces = true
		}
		if source.RemoveTrailingDot == nil {
			config.Sources.Terms[i].RemoveTrailingDot = new(bool)
			*config.Sources.Terms[i].RemoveTrailingDot = true
		}
		if source.ItemsDelimiter == "" {
			config.Sources.Terms[i].ItemsDelimiter = ","
		}
	}

	for i, source := range config.Sources.Dates {
		if source.TrimSpaces == nil {
			config.Sources.Dates[i].TrimSpaces = new(bool)
			*config.Sources.Dates[i].TrimSpaces = true
		}
		if source.RemoveTrailingDot == nil {
			config.Sources.Dates[i].RemoveTrailingDot = new(bool)
			*config.Sources.Dates[i].RemoveTrailingDot = true
		}
		if source.ItemsDelimiter == "" {
			config.Sources.Dates[i].ItemsDelimiter = ","
		}
	}

	if config.Exporter.DeleteTypstFiles == nil {
		config.Exporter.DeleteTypstFiles = new(bool)
		*config.Exporter.DeleteTypstFiles = true
	}
}
