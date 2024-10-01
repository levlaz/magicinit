//A Dagger pipeline for go generated via Magicinit™
//                        _      _       _ _
//                       (_)    (_)     (_) |
//  _ __ ___   __ _  __ _ _  ___ _ _ __  _| |_
// | '_ ` _ \ / _` |/ _` | |/ __| | '_ \| | __|
// | | | | | | (_| | (_| | | (__| | | | | | |_
// |_| |_| |_|\__,_|\__, |_|\___|_|_| |_|_|\__|
//                   __/ |
//                  |___/

package main

import (
	"context"
	"dagger/go/internal/dagger"
	"fmt"
)

type Go struct {
	GoVersion string
	Source    *dagger.Directory
}

func New(
	goVersion string,
	// +defaultPath="/"
	source *dagger.Directory,
) *Go {
	return &Go{
		GoVersion: goVersion,
		Source:    source,
	}
}

func (m *Go) Base(ctx context.Context) *dagger.Container {
	return dag.
		Container().
		From(fmt.Sprintf("golang:%s", m.GoVersion)).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src")
}

func (m *Go) Build(ctx context.Context) *dagger.Container {
	return m.
		Base(ctx).
		WithExec([]string{"go", "build", "-o", "main", "."})
}

func (m *Go) Lint(ctx context.Context) *dagger.Container {
	return m.
		Base(ctx).
		WithExec([]string{"go", "vet", "."})
}

func (m *Go) Test(ctx context.Context) *dagger.Container {
	return m.
		Base(ctx).
		WithExec([]string{"go", "test", "-v", "."})
}
