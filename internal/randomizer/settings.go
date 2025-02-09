package randomizer

import (
	"encoding/json"
	"io"
	"os"
)

type GeneratorSettings struct {
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
	return json.Marshal(g)
}
