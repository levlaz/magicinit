package main

import (
	"context"
	"dagger/magicinit/internal/dagger"
	"fmt"
)


type SourceInspect struct {
	Language string
	Version string
}


type Inspector interface {
	// Lookup returns true if the source directory matches the inspector language
	Lookup(context.Context, *dagger.Directory) (bool, error)

	// Inspects the source directory and returns the language and version
	//
	// This function should be called after Lookup.
	Inspect(context.Context, *dagger.Directory) (*SourceInspect, error)
}

func List() map[string]Inspector {
	return map[string]Inspector{
		"go": &goInspector{},
		"python": &pythonInspector{},
		"ruby": &rubyInspector{},
		"typescript": &typescriptInspector{},
	}
}

func Get(language string) (Inspector, error) {
	inspector, ok := List()[language]
	if !ok {
		return nil, fmt.Errorf("inspector not found for language %s", language)
	}

	return inspector, nil
} 