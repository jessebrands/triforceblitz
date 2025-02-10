package randomizer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

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
		"--no_log",
	)
	// Create the log file.
	logfile, err := os.Create(filepath.Join(t.settings.OutputDir, "TriforceBlitz.log"))
	if err != nil {
		return err
	}
	defer logfile.Close()
	// The randomizer only writes to stderr, so grab the stderr pipe.
	// We'll use this to forward log messages to the caller.
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer stderr.Close()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			text := strings.TrimRight(scanner.Text(), "\r\n\t ")
			// Old versions of the randomizer display this annoying, bugged header that isn't
			// escaped properly and dumps a ton of Python source code out into the log.
			// We do not like that so we kill it here.
			if strings.Contains(text, ", and the latest is version ") {
				continue
			}
			if text != "" {
				_, _ = fmt.Fprintf(logfile, "[%v]  %s\n",
					time.Now().Format(time.StampMilli),
					text)
			}
			if t.OnMessage != nil {
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
	// Ensure output directory exists.
	if err := os.MkdirAll(t.settings.OutputDir, 0755); err != nil {
		return err
	}
	if err := t.settings.WriteFile(settingsFile); err != nil {
		return err
	}
	if err := t.run(interpreter, settingsFile); err != nil {
		return err
	}
	return t.verifyGeneratedFiles()
}
