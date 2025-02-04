package generator_test

import (
	"log"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/jessebrands/triforceblitz/internal/generator"
)

var fsys = fstest.MapFS{
	"v1.0.0-blitz-1.0/Generator.json": {
		Mode: 0644,
		Data: []byte(`
			{
				"version": "1.0.0-blitz-1.0",
				"presets": [
					{
						"id": "default",
						"preset": "Triforce Blitz"
					}
				]
			}
		`),
	},
}

func assertGenerator(t *testing.T, got *generator.Generator, want *generator.Generator) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\n\tgot  %+v\n\twant %+v", got, want)
	}
}

func TestFindGenerators(t *testing.T) {
	generators, err := generator.FindGeneratorsFromFS(fsys, ".")
	if err != nil {
		t.Fatal(err)
	}
	if len(generators) != 1 {
		t.Errorf("got %d generator paths, wanted 1", len(generators))
	}
	got := generators[0]
	want := &generator.Generator{
		Path:    "v1.0.0-blitz-1.0",
		Version: "1.0.0-blitz-1.0",
		Presets: []generator.Preset{
			{
				Id:     "default",
				Preset: "Triforce Blitz",
			},
		},
	}
	assertGenerator(t, got, want)
}

func init() {
	err := fstest.TestFS(fsys,
		"v1.0.0-blitz-1.0/Generator.json",
	)
	if err != nil {
		log.Fatal("test filesystem failed test")
	}
}
