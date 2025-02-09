package randomizer

const (
	DefaultPreset = "default"
)

type Preset struct {
	Value   string `json:"preset"`
	Ordinal int    `json:"ordinal,omitempty"`
}

type PresetMap map[string]Preset

func (m *PresetMap) Default() (Preset, error) {
	if p, ok := (*m)[DefaultPreset]; ok {
		return p, nil
	} else {
		return Preset{}, ErrNoDefaultPreset
	}
}
