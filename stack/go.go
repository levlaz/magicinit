package stack

import (
	"bytes"
	"context"
	"dagger/magicinit/inspection"
	"dagger/magicinit/internal/dagger"
	"fmt"
	"log"
	"regexp"
	"slices"
	"text/template"
)

type goStack struct{}

func (g *goStack) Lookup(ctx context.Context, source *dagger.Directory) (bool, error) {
	files, err := source.Entries(ctx)
	if err != nil {
		return false, fmt.Errorf("Go.Lookup: failed to list source directory entries: %w", err)
	}

	if slices.Contains(files, "go.mod") {
		return true, nil
	}

	return false, nil
}

func (g *goStack) Inspect(ctx context.Context, source *dagger.Directory) (*inspection.Source, error) {
	var inspection = &inspection.Source{
		Language: "go",
	}

	goModContent, err := source.File("go.mod").Contents(ctx)
	if err != nil {
		log.Printf("Go.Inspect: failed to read go.mod: %v", err)
		return inspection, nil
	}

	re := regexp.MustCompile(`go\s+([\d.]+)`)
	matches := re.FindStringSubmatch(goModContent)
	if len(matches) > 1 {
		inspection.Version = matches[1]
	}

	return inspection, nil	
}

func (g *goStack) Init(ctx context.Context, inspection *inspection.Source) (*dagger.Directory, error) {
	goTemplateDir := 	dagger.Connect().CurrentModule().Source().Directory("templates/go")
	goTemplate, err := goTemplateDir.File("main.go.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initGo: failed to read go template: %w", err)
	}

	var goTemplateData = struct {
		GoVersion string
		Compose   bool
	}{
		GoVersion: inspection.Version,
		Compose:   inspection.Compose,
	}

	tmpl, err := template.New("go.tmpl").Parse(goTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initGo: failed to parse go template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, goTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initGo: failed to execute go template: %w", err)
	}

	return goTemplateDir.WithNewFile("main.go", buf.String()).WithoutFile("main.go.tmpl"), nil

}