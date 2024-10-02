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

	"github.com/pelletier/go-toml/v2"
)

var pythonFiles = []string{
	"pyproject.toml",
	"setup.py",
	"requirements.txt",
}

type pythonStack struct{}

func (p *pythonStack) Lookup(ctx context.Context, source *dagger.Directory) (bool, error) {
	files, err := source.Entries(ctx)
	if err != nil {
		return false, fmt.Errorf("Python.Lookup: failed to list source directory entries: %w", err)
	}

	isPython := slices.ContainsFunc(files, func(file string) bool {
		return slices.Contains(pythonFiles, file)
	})

	return isPython, nil
}

func (p *pythonStack) Inspect(ctx context.Context, source *dagger.Directory) (*inspection.Source, error) {
	var pyProjectToml struct {
		Project struct {
			RequiresPython string `toml:"requires-python"`
		}
	}

	var inspection = &inspection.Source{
		Language: "python",
	}

	pyprojectContent, err := source.File("pyproject.toml").Contents(ctx)
	if err != nil {
		log.Printf("Python.Inspect: failed to read pyproject.toml: %v", err)
		return inspection, nil
	}

	if err := toml.Unmarshal([]byte(pyprojectContent), &pyProjectToml); err != nil {
		log.Printf("Python.Inspect: failed to unmarshal pyproject.toml: %v", err)
		return inspection, nil
	}

	version := pyProjectToml.Project.RequiresPython
	re := regexp.MustCompile(`[><=~!]*\s*([\d.]+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) > 1 {
		inspection.Version = matches[1]
	}

	return inspection, nil
}

func (p *pythonStack) Init(ctx context.Context, inspection *inspection.Source) (*dagger.Directory, error) {
	pythonTemplateDir := dagger.Connect().CurrentModule().Source().Directory("templates/python")
	pythonTemplate, err := pythonTemplateDir.File("/src/main/main.py.tmpl").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initPython: failed to read python template: %w", err)
	}

	var pythonTemplateData = struct {
		PythonVersion string
		Compose       bool
	}{
		PythonVersion: inspection.Version,
		Compose:       inspection.Compose,
	}

	tmpl, err := template.New("python.tmpl").Parse(pythonTemplate)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initPython: failed to parse python template: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pythonTemplateData)
	if err != nil {
		return nil, fmt.Errorf("Magicinit.initPython: failed to execute python template: %w", err)
	}

	return pythonTemplateDir.WithNewFile("/src/main/main.py", buf.String()).WithoutFile("/src/main/main.py.tmpl"), nil

}