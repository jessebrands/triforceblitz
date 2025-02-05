package main

import (
	"context"
)

type Source interface {
	// ListAvailable returns a list of available Packages from this Source.
	ListAvailable(ctx context.Context) ([]Package, error)

	// Type returns a string identifying the type of the Source.
	Type() string

	// Name returns a name uniquely identifying this Source.
	Name() string
}

func SourceIdentifier(source Source) string {
	return source.Type() + ":" + source.Name()
}
