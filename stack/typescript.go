package stack

import (
	"bytes"
	"context"
	"dagger/magicinit/inspection"
	"dagger/magicinit/internal/dagger"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"slices"
	"text/template"
)

type typescriptStack struct{}

func (t *typescriptStack) Lookup(ctx context.Context, source *dagger.Directory) (bool, error) {
	files, err := source.Entries(ctx)
	if err != nil {
		return false, fmt.Errorf("Typescript.Lookup: failed to list source directory entries: %w", err)
	}

	if slices.Contains(files, "package.json") {
		return true, nil
	}

	return false, nil
}

func (t *typescriptStack) Inspect(ctx context.Context, source *dagger.Directory) (*inspection.Source, error) {
	var packageJson struct {
		Engines struct {
			Node string `json:"node"`
		}
	}

	var inspection = &inspection.Source{
		Language: "typescript",
	}

	packageJsonContent, err := source.File("package.json").Contents(ctx)

	if err != nil {
		log.Printf("Typescript.Inspect: failed to read package.json: %v", err)
		return inspection, nil
	}

	if err := json.Unmarshal([]byte(packageJsonContent), &packageJson); err != nil {
		log.Printf("Typescript.Inspect: failed to unmarshal package.json: %v", err)
		return inspection, nil
	}

	version := packageJson.Engines.Node
	re := regexp.MustCompile(`[><=~!]*\s*([\d.]+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) > 1 {
		inspection.Version = matches[1]
	}

	return inspection, nil
}

func (t *typescriptStack) Init(ctx context.Context, inspection *inspection.Source) (*dagger.Directory, error) {
	typescriptTemplateDir := dagger.Connect().CurrentModule().Source().Directory("templates/typescript")
	typescriptTemplate, err := typescriptTemplateDir.File("src/index.ts.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to read typescript template: %w", err)
	}

	daggerJsonTemplate, err := typescriptTemplateDir.File("dagger.json.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to read typescript dagger.json template: %w", err)
	}

	var typescriptTemplateData = struct {
		TypescriptVersion string
		Compose           bool
	}{
		TypescriptVersion: inspection.Version,
		Compose:           inspection.Compose,
	}

	srcTmpl, err := template.New("typescript.tmpl").Parse(typescriptTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to parse typescript template: %w", err)
	}

	var buf bytes.Buffer
	err = srcTmpl.Execute(&buf, typescriptTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to execute typescript template: %w", err)
	}

	daggerJsonTmpl, err := template.New("dagger-json.tmpl").Parse(daggerJsonTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to read typescript dagger.json template: %w", err)
	}

	var daggerJsonBuf bytes.Buffer
	err = daggerJsonTmpl.Execute(&daggerJsonBuf, typescriptTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initTypescript: failed to execute typescript dagger.json template: %w", err)
	}

	return typescriptTemplateDir.WithNewFile("src/index.ts", buf.String()).
		WithoutFile("src/index.ts.tmpl").
		WithNewFile("dagger.json", daggerJsonBuf.String()).
		WithoutFile("dagger.json.tmpl"), nil
}