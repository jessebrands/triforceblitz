package randomizer

import (
	"os/exec"
	"path/filepath"

	"github.com/jessebrands/triforceblitz/internal/python"
)

const (
	EntrypointFilename = "OoTRandomizer.py"
	SettingsFilename   = "settings.json"
)

const (
	SettingsPresetArg = "--settings_preset"
	SettingsFileArg   = "--settings"
)

type Generator struct {
	Version Version
	Path    string
	Presets PresetMap
}

type GenerateSeedOpts struct {
	Seed      string
	Preset    string
	OutputDir string
	RomFile   string
}

func (g *Generator) Generate(interpreter python.Interpreter, opts GenerateSeedOpts) (*exec.Cmd, error) {
	// Create the settings file first.
	settings := NewSettings(opts.Seed, opts.OutputDir, opts.RomFile)
	settingsFile := filepath.Join(settings.OutputDir, SettingsFilename)
	if err := settings.WriteFile(settingsFile); err != nil {
		return nil, err
	}
	// Call the generator, plain and simple!
	cmd := interpreter.Command(
		g.Entrypoint(),
		SettingsPresetArg, opts.Preset,
		SettingsFileArg, settingsFile,
	)
	return cmd, nil
}

func (g *Generator) Entrypoint() string {
	return filepath.Join(g.Path, EntrypointFilename)
}

func (g *Generator) String() string {
	return g.Version.String()
}

func (p *Preset) String() string {
	return p.Value
}
