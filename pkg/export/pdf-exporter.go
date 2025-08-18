package export

import (
	"archivist/pkg/pipeline"
	_ "embed"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed document.tmpl
var documentTemplate string

type PdfExporter struct {
	parceling        Parcel
	deleteTypstFiles bool
}

func CreatePdfExporter(personsParts, termsParts, datesParts int, deleteTypstFiles bool) Exporter {
	return &PdfExporter{
		parceling: Parcel{
			PersonsParts: personsParts,
			TermsParts:   termsParts,
			DatesParts:   datesParts,
		},
		deleteTypstFiles: deleteTypstFiles,
	}
}

func (p *PdfExporter) Export(data pipeline.Result) error {
	return p.ExportWithParceling(data)
}

func (p *PdfExporter) ExportWithParceling(data pipeline.Result) error {
	persons := groupInOrder(data.Persons)
	terms := groupInOrder(data.Terms)
	dates := groupInOrder(data.Dates)

	for _, person := range persons {
		parcels := parcelTagGroup(person, p.parceling.PersonsParts)
		for i, parcel := range parcels {
			err := RenderPdfWithTypstParceled(parcel, i, len(parcels), p.deleteTypstFiles)
			if err != nil {
				return err
			}
		}
	}

	for _, term := range terms {
		parcels := parcelTagGroup(term, p.parceling.TermsParts)
		for i, parcel := range parcels {
			err := RenderPdfWithTypstParceled(parcel, i, len(parcels), p.deleteTypstFiles)
			if err != nil {
				return err
			}
		}
	}

	for _, date := range dates {
		parcels := parcelTagGroup(date, p.parceling.DatesParts)
		for i, parcel := range parcels {
			err := RenderPdfWithTypstParceled(parcel, i, len(parcels), p.deleteTypstFiles)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func parcelTagGroup(group TagGroup, partsAmount int) []TagGroup {
	if partsAmount <= 1 {
		return []TagGroup{group}
	}

	var totalItems int
	for _, topic := range group.Topics {
		totalItems += len(topic.Items)
	}

	if totalItems <= partsAmount {
		return []TagGroup{group}
	}

	maxItems := int(math.Ceil(float64(totalItems) / float64(partsAmount)))

	var parcels []TagGroup
	currentParcel := TagGroup{Tag: group.Tag, Topics: []TopicGroup{}}
	var currentItems int

	for _, topic := range group.Topics {
		for _, item := range topic.Items {
			if currentItems >= maxItems {
				parcels = append(parcels, currentParcel)
				currentParcel = TagGroup{Tag: group.Tag, Topics: []TopicGroup{}}
				currentItems = 0
			}

			var targetTopic *TopicGroup
			for i := range currentParcel.Topics {
				if currentParcel.Topics[i].Topic == topic.Topic {
					targetTopic = &currentParcel.Topics[i]
					break
				}
			}

			if targetTopic == nil {
				currentParcel.Topics = append(currentParcel.Topics, TopicGroup{
					Topic: topic.Topic,
					Items: []ShortItem{},
				})
				targetTopic = &currentParcel.Topics[len(currentParcel.Topics)-1]
			}

			targetTopic.Items = append(targetTopic.Items, item)
			currentItems++
		}
	}

	if currentItems > 0 {
		parcels = append(parcels, currentParcel)
	}

	return parcels
}

func RenderPdfWithTypstParceled(group TagGroup, parcelIndex, totalParcels int, deleteTypstFile bool) error {
	EnsureDataFolder()

	filename := "output/" + group.Tag
	if totalParcels > 1 {
		filename += fmt.Sprintf("_%d", parcelIndex+1)
	}
	filename += ".typ"

	filenamePdf := strings.TrimSuffix(filename, ".typ") + ".pdf"

	f, err := os.Create(filename)
	defer func() { _ = f.Close() }()
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"highlight": HighlightItem,
	}

	t := template.Must(template.New("doc").
		Funcs(funcMap).
		Parse(documentTemplate))

	err = t.Execute(f, group)
	if err != nil {
		return err
	}

	cmd := exec.Command("typst", "compile", filename)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to compile file %s:\n%s", filename, string(out))
	}

	log.Printf("Created file %s", filenamePdf)

	if deleteTypstFile {
		_ = os.Remove(filename)
	}

	return nil
}

func HighlightItem(summary string) string {
	bracketIndex := strings.Index(summary, "(")
	dashIndex := strings.Index(summary, "—")
	if bracketIndex == -1 && dashIndex == -1 {
		return summary
	}
	if bracketIndex != -1 && dashIndex != -1 {
		m := min(dashIndex, bracketIndex)
		return "*" + summary[:m] + "*" + summary[m:]
	}
	if bracketIndex == -1 {
		return "*" + summary[:dashIndex] + "*" + summary[dashIndex:]
	} else {
		return "*" + summary[:bracketIndex] + "*" + summary[bracketIndex:]
	}
}

func EnsureDataFolder() {
	const path = "output"

	if _, e := os.Stat(path); os.IsNotExist(e) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Fatalf("Failed to create directory %s: %s", path, err)
		}
		log.Printf("Created data directory: %s", path)
	}
}
