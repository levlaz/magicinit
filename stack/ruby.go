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

type rubyStack struct{}

func (r *rubyStack) Lookup(ctx context.Context, source *dagger.Directory) (bool, error) {
	files, err := source.Entries(ctx)
	if err != nil {
		return false, fmt.Errorf("Ruby.Lookup: failed to list source directory entries: %w", err)
	}

	if slices.Contains(files, "Gemfile") {
		return true, nil
	}

	return false, nil
}

func (r *rubyStack) Inspect(ctx context.Context, source *dagger.Directory) (*inspection.Source, error) {
	var inspection = &inspection.Source{
		Language: "ruby",
	}

	gemfileContent, err := source.File("Gemfile").Contents(ctx)
	if err != nil {
		log.Printf("Ruby.Inspect: failed to read Gemfile: %v", err)
		return inspection, nil
	}

	re := regexp.MustCompile(`ruby\s+['"][><=~!]+\s*([\d.]+)['"]`)
	matches := re.FindStringSubmatch(gemfileContent)
	if len(matches) > 1 {
		inspection.Version = matches[1]
	}

	return inspection, nil
}

func (r *rubyStack) Init(ctx context.Context, inspection *inspection.Source) (*dagger.Directory, error) {
	rubyTemplateDir := dagger.Connect().CurrentModule().Source().Directory("templates/ruby")
	rubyTemplate, err := rubyTemplateDir.File("src/index.ts.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initRuby: failed to read ruby template: %w", err)
	}

	var rubyTemplateData = struct {
		RubyVersion string
		Compose     bool
	}{
		RubyVersion: inspection.Version,
		Compose:     inspection.Compose,
	}

	tmpl, err := template.New("ruby.tmpl").Parse(rubyTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initRuby: failed to parse ruby template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, rubyTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initRuby: failed to execute ruby template: %w", err)
	}

	return rubyTemplateDir.WithNewFile("src/index.ts", buf.String()).WithoutFile("src/index.ts.tmpl"), nil

}