package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jessebrands/triforceblitz/internal/python"
	"github.com/jessebrands/triforceblitz/internal/randomizer"
)

var (
	svc = randomizer.NewService(os.Getenv("TRIFORCEBLITZ_GENERATORS_DIR"))
)

type GenerateSeedOpts struct {
	Seed    string
	Version randomizer.Version
	RomFile string
	OutDir  string
}

func ParseGenerateSeedOpts(args []string) (GenerateSeedOpts, error) {
	opts := GenerateSeedOpts{}
	flags := flag.NewFlagSet("generate", flag.ExitOnError)
	flags.StringVar(&opts.RomFile, "R", "", "ROM file to use")
	flags.StringVar(&opts.OutDir, "o", "", "directory to store the result in")
	flags.StringVar(&opts.Seed, "s", "", "random number generator seed passed to the generator")
	flags.Var(&opts.Version, "r", "generator version to use")
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
	// Generate the seed.
	generatorOpts := randomizer.GenerateSeedOpts{
		OutputDir: opts.OutDir,
		Seed:      opts.Seed,
		RomFile:   opts.RomFile,
		Preset:    "Triforce Blitz",
	}
	fmt.Printf("Generating seed %s with generator %s\n", opts.Seed, opts.Version.String())
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
		fmt.Printf("Generator process error: %s\n", err.Error())
		os.Exit(1)
	}
	if err := cmd.Wait(); err != nil {
		fmt.Printf("Wait(): %s\n", err.Error())
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
