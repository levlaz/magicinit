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

	"golang.org/x/sync/errgroup"
)

type Go struct {
	GoVersion string
	Source    *dagger.Directory
}

func New(
	// +default="1.23.0"
	goVersion string,

	// +defaultPath="/"
	source *dagger.Directory,
) *Go {
	return &Go{
		GoVersion: goVersion,
		Source:    source,
	}
}

// Base returns a container with the go project mounted
func (m *Go) Base() *dagger.Container {
	return dag.
		Container().
		From(fmt.Sprintf("golang:%s", m.GoVersion)).
		WithMountedDirectory("/src", m.Source).
		WithWorkdir("/src")
}

// Build returns a container with the go project mounted and builds the project
func (m *Go) Build() *dagger.Container {
	return m.
		Base().
		WithExec([]string{"go", "build", "-o", "main", "."})
}

// Lint returns a container with the go project mounted and runs the go vet command
func (m *Go) Lint() *dagger.Container {
	return m.
		Base().
		WithExec([]string{"go", "vet", "."})
}

// Test returns a container with the go project mounted and runs the go test command
func (m *Go) Test() *dagger.Container {
	return m.
		Base().
		WithExec([]string{"go", "test", "-v", "."})
}

// Run entire pipleine
//
// uses error group to run all commands in parallel
func (m *Go) Ci(ctx context.Context) error {

	eg, gctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		_, err := m.Build().Sync(gctx)
		return err
	})

	eg.Go(func() error {
		_, err := m.Lint().Sync(gctx)
		return err
	})

	eg.Go(func() error {
		_, err := m.Test().Sync(gctx)
		return err
	})

	return eg.Wait()
}
