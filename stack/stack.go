package stack

import (
	"context"
	"dagger/magicinit/inspection"
	"dagger/magicinit/internal/dagger"
	"fmt"
)


type Stack interface {
	// Lookup returns true if the source directory matches the stack language
	Lookup(context.Context, *dagger.Directory) (bool, error)

	// Inspects the source directory and returns the language and version
	//
	// This function should be called after Lookup.
	Inspect(context.Context, *dagger.Directory) (*inspection.Source, error)

	Init(ctx context.Context, inspection *inspection.Source) (*dagger.Directory, error)
}

func List() map[string]Stack {
	return map[string]Stack{
		"go": &goStack{},
		"python": &pythonStack{},
		"ruby": &rubyStack{},
		"typescript": &typescriptStack{},
	}
}

func Get(language string) (Stack, error) {
	inspector, ok := List()[language]
	if !ok {
		return nil, fmt.Errorf("inspector not found for language %s", language)
	}

	return inspector, nil
} 