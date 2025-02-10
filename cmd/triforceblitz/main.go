package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/jessebrands/triforceblitz/internal/config"
	"os"
	"strings"
	"time"

	"github.com/jessebrands/triforceblitz/internal/python"
	"github.com/jessebrands/triforceblitz/internal/randomizer"
)

var (
	svc = randomizer.NewService(config.GetGeneratorDir())
)

type GenerateSeedOpts struct {
	Seed    string
	Preset  string
	Version randomizer.Version
	RomFile string
	OutDir  string
}

func ParseGenerateSeedOpts(args []string) (GenerateSeedOpts, error) {
	opts := GenerateSeedOpts{}
	flags := flag.NewFlagSet("generate", flag.ExitOnError)
	flags.StringVar(&opts.Seed, "s", "", "random number generator seed passed to the generator")
	flags.StringVar(&opts.Preset, "p", "default", "settings preset to use")
	flags.Var(&opts.Version, "r", "generator version to use")
	flags.StringVar(&opts.RomFile, "R", "", "ROM file to use")
	flags.StringVar(&opts.OutDir, "o", "", "directory to store the result in")
	if err := flags.Parse(args); err != nil {
		return opts, err
	}
	if opts.RomFile == "" {
		return opts, errors.New("no ROM file specified")
	}
	if opts.Seed == "" {
		if seed, err := randomizer.GenerateSeedString(32); err != nil {
			return opts, err
		} else {
			opts.Seed = seed
		}
	}
	if opts.OutDir == "" {
		if dir, err := os.Getwd(); err != nil {
			return opts, err
		} else {
			opts.OutDir = dir
		}
	}
	return opts, nil
}

func selectPreset(generator *randomizer.Generator, preset string) (string, error) {
	// Empty string means we should just grab the default preset.
	if preset == "" {
		if p, err := generator.Presets.Default(); err == nil {
			return p.Value, nil
		} else {
			return "", err
		}
	}
	// A preset was specified, this might be a 'symbolic' preset.
	// Look it up in our map and return it if we find it.
	if p, ok := generator.Presets[preset]; ok {
		return p.Value, nil
	}
	// The preset was not found symbolically but was set.
	// This might be a preset literal, so we will just return it as-is.
	return preset, nil
}

func generateSeed(args []string) {
	opts, err := ParseGenerateSeedOpts(args)
	if err != nil {
		fmt.Printf("Could not generate seed: %s\n", err.Error())
		os.Exit(1)
	}
	// Check if the actual generator even exists.
	generator, err := svc.GetGenerator(opts.Version)
	if err != nil {
		fmt.Printf("Could not use generator %s: %s.\n", opts.Version.String(), err.Error())
		os.Exit(1)
	}
	// Find a Python interpreter.
	interpreter, err := python.FindInterpreter()
	if err != nil {
		fmt.Printf("Could not find Python interpreter: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Found Python interpreter: %s\n", interpreter.Path())
	// Select the preset.
	preset, err := selectPreset(generator, opts.Preset)
	if err != nil {
		fmt.Printf("Could not select preset: %s\n", err.Error())
		os.Exit(1)
	}
	// Generate the seed.
	generatorOpts := randomizer.GenerateSeedOpts{
		OutputDir: opts.OutDir,
		Seed:      opts.Seed,
		RomFile:   opts.RomFile,
		Preset:    preset,
	}
	fmt.Printf("Generating seed %s with generator %s using settings preset %s\n",
		opts.Seed,
		opts.Version.String(),
		preset)
	start := time.Now()
	cmd, err := generator.Generate(interpreter, generatorOpts)
	if err != nil {
		fmt.Printf("Failed to invoke generator: %s\n", err.Error())
		os.Exit(1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Failed to get stderr from generator: %s\n", err.Error())
		os.Exit(1)
	}
	defer stderr.Close()
	scanner := bufio.NewScanner(stderr)
	// At this point, we need to figure out a way to grab the output...
	go func() {
		for scanner.Scan() {
			line := strings.TrimRight(scanner.Text(), " \t\r\n")
			fmt.Printf("==> %s\n", line)
		}
	}()
	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start randomizer: %s\n", err.Error())
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Failed to generate seed because of randomizer error: %s\n", err.Error())
		os.Exit(1)
	}
	elapsed := time.Since(start)
	fmt.Printf("Seed generation finished in %s\n", elapsed)
}

func main() {
	command := os.Args[1]
	switch command {
	case "generate":
		generateSeed(os.Args[2:])

	default:
		fmt.Println("You must specify a command, type 'triforceblitz help' for a list of commands.")
	}
}

func init() {
	if err := svc.Synchronize(); err != nil {
		panic(err)
	}
}
