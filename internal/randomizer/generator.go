package randomizer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jessebrands/triforceblitz/internal/python"
)

const (
	EntrypointFilename = "OoTRandomizer.py"
	SettingsFilename   = "settings.json"
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

type GenerateSeedOpts struct {
	Seed      string
	Preset    string
	OutputDir string
	RomFile   string
}

type generatorSettings struct {
	// String used to seed the randomizer.
	Seed string `json:"seed"`

	// Where to store generated files.
	OutputDir string `json:"output_dir"`

	// Filename of outputted files.
	OutputFilename string `json:"output_file"`

	// Path to the ROM file.
	RomFile string `json:"rom"`

	// Whether to create a patch file.
	CreatePatch bool `json:"create_patch_file"`

	// Whether to create a compressed ROM file.
	CreateCompressedRom bool `json:"create_compressed_rom"`

	// Whether to output a cosmetics log.
	CreateCosmeticsLog bool `json:"create_cosmetics_log"`

	// Legacy setting that controls ROM output.
	CompressRom string `json:"compress_rom"`
}

func createSettingsFile(settings generatorSettings, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

func defaultSettings(randomizerSeed string, outDir string, romFile string) generatorSettings {
	return generatorSettings{
		Seed:                randomizerSeed,
		OutputDir:           outDir,
		RomFile:             romFile,
		OutputFilename:      "TriforceBlitz",
		CreatePatch:         true,
		CreateCompressedRom: false,
		CreateCosmeticsLog:  false,
		CompressRom:         "Patch",
	}
}

func (g *Generator) Generate(interpreter python.Interpreter, opts GenerateSeedOpts) error {
	// Create the settings file first.
	settings := defaultSettings(opts.Seed, opts.OutputDir, opts.RomFile)
	settingsFile := filepath.Join(settings.OutputDir, SettingsFilename)
	if err := createSettingsFile(settings, settingsFile); err != nil {
		return err
	}
	// Call the generator, plain and simple!
	cmd := interpreter.Command(
		g.Entrypoint(),
		"--settings_preset", opts.Preset,
		"--settings", settingsFile,
	)
	out, err := cmd.CombinedOutput()
	fmt.Print(string(out))
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) Entrypoint() string {
	return filepath.Join(g.Path, EntrypointFilename)
}

func (g *Generator) String() string {
	return g.Version.String()
}

func (p *Preset) String() string {
	return p.Preset
}
