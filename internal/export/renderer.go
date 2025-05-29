package export

import (
	"bismark/internal/pipeline"
	_ "embed"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed document.tmpl
var documentTemplate string

func ToPdf(data pipeline.Result) error {
	persons := groupInOrder(data.Persons)
	terms := groupInOrder(data.Terms)
	dates := groupInOrder(data.Dates)
	for _, person := range persons {
		err := RenderFile(person, "")
		if err != nil {
			return err
		}
	}

	for _, term := range terms {
		err := RenderFile(term, "")
		if err != nil {
			return err
		}
	}

	for _, date := range dates {
		err := RenderFile(date, "")
		if err != nil {
			return err
		}
	}

	return nil
}

func RenderFile(group TagGroup, filePostfix string) error {
	filename := group.Tag + filePostfix + ".typ"
	filenamePdf := group.Tag + filePostfix + ".pdf"

	f, err := os.Create(filename)
	defer f.Close()
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
