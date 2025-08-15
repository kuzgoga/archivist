package export

import "archivist/pkg/pipeline"

type Exporter interface {
	Export(data pipeline.Result) error
}
