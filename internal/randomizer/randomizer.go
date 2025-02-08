package randomizer

const (
	EntrypointFilename = "OoTRandomizer.py"
	DefaultPreset      = "default"
)

type Randomizer struct {
	Version Version
	Path    string
	Presets []Preset
}

type Preset struct {
	Id      string `json:"id"`
	Preset  string `json:"preset"`
	Ordinal int    `json:"ordinal,omitempty"`
}

func (g *Randomizer) String() string {
	return g.Version.String()
}

func (p *Preset) String() string {
	return p.Preset
}
