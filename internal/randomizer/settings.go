package randomizer

import (
	"encoding/json"
	"io"
	"os"
)

type GeneratorSettings struct {
	// String used to seed the randomizer.
	Seed string

	// Where to store generated files.
	OutputDir string

	// Filename of outputted files.
	OutputFilename string

	// Path to the ROM file.
	RomFile string

	// Whether to create a patch file.
	CreatePatch bool

	// Whether to create a compressed ROM file.
	CreateCompressedRom bool

	// Whether to output a cosmetics log.
	CreateCosmeticsLog bool

	// Legacy setting that controls ROM output.
	CompressRom string
}

func NewSettings(randomizerSeed string, outDir string, romFile string) GeneratorSettings {
	return GeneratorSettings{
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

func (g *GeneratorSettings) WriteFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	return g.WriteJSON(f)
}

func (g *GeneratorSettings) WriteJSON(w io.Writer) error {
	b, err := g.MarshalJSON()
	if err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func (g *GeneratorSettings) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Seed                string `json:"seed"`
		OutputDir           string `json:"output_dir"`
		OutputFilename      string `json:"output_file"`
		RomFile             string `json:"rom"`
		CreatePatch         bool   `json:"create_patch_file"`
		CreateCompressedRom bool   `json:"create_compressed_rom"`
		CreateCosmeticsLog  bool   `json:"create_cosmetics_log"`
		CompressRom         string `json:"compress_rom"`
	}{
		Seed:                g.Seed,
		OutputDir:           g.OutputDir,
		OutputFilename:      g.OutputFilename,
		RomFile:             g.RomFile,
		CreatePatch:         g.CreatePatch,
		CreateCompressedRom: g.CreateCompressedRom,
		CreateCosmeticsLog:  g.CreateCosmeticsLog,
		CompressRom:         g.CompressRom,
	})
}
