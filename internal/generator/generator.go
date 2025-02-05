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
	Id      string `json:"id"`
	Preset  string `json:"preset"`
	Ordinal int    `json:"ordinal,omitempty"`
}

func (g *Generator) String() string {
	return g.Version.String()
}

func (p *Preset) String() string {
	return p.Preset
}
