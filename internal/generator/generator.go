package generator

const (
	MetadataFilename   = ".generator-metadata.json"
	EntrypointFilename = "OoTRandomizer.py"
	DefaultPreset      = "default"
)

type Generator struct {
	Version Version
	Path    string
	Presets []Preset
}

type Preset struct {
	Id      string
	Preset  string
	Ordinal int
}

func (g *Generator) String() string {
	return g.Version.String()
}

func (p *Preset) String() string {
	return p.Preset
}
