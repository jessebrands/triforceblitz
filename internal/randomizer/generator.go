package randomizer

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

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

func (g *Generator) CreateTask(seed, preset, outputDir, romFile string) *GeneratorTask {
	return &GeneratorTask{
		generator: g,
		preset:    preset,
		settings:  NewSettings(seed, outputDir, romFile),
	}
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

type GeneratorTask struct {
	generator *Generator
	preset    string
	settings  GeneratorSettings

	// OnMessage is called while the randomizer runs to report its
	// messages.
	//
	// This callback is called from a goroutine.
	OnMessage func(string)
}

func (t *GeneratorTask) run(interpreter python.Interpreter, settingsFile string) error {
	cmd := interpreter.Command(
		t.generator.Entrypoint(),
		SettingsPresetArg, t.preset,
		SettingsFileArg, settingsFile,
	)
	// The randomizer only writes to stderr, so grab the stderr pipe.
	// We'll use this to forward log messages to the caller.
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()
	go func() {
		if t.OnMessage != nil {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				text := strings.TrimRight(scanner.Text(), "\r\n\t ")
				t.OnMessage(text)
			}
		}
	}()
	// We're all set up, let's run the command.
	return cmd.Run()
}

func (t *GeneratorTask) verifyGeneratedFiles() error {
	if t.settings.CreatePatch {
		patchFile := filepath.Join(t.settings.OutputDir, t.settings.OutputFilename+".zpf")
		if _, err := os.Stat(patchFile); os.IsNotExist(err) {
			return ErrNoPatchFile
		}
	}
	spoilerFile := filepath.Join(t.settings.OutputDir, t.settings.OutputFilename+"_Spoiler.json")
	if _, err := os.Stat(spoilerFile); os.IsNotExist(err) {
		return ErrNoSpoilerLog
	}
	return nil
}

func (t *GeneratorTask) Generate(interpreter python.Interpreter) error {
	settingsFile := filepath.Join(t.settings.OutputDir, SettingsFilename)
	if err := t.settings.WriteFile(settingsFile); err != nil {
		return err
	}
	if err := t.run(interpreter, settingsFile); err != nil {
		return err
	}
	return t.verifyGeneratedFiles()
}
