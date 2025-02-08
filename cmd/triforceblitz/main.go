package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jessebrands/triforceblitz/internal/python"
	"github.com/jessebrands/triforceblitz/internal/randomizer"
	"github.com/jessebrands/triforceblitz/internal/seed"
)

var (
	svc = randomizer.NewService(os.Getenv("TRIFORCEBLITZ_RANDOMIZERS_DIR"))
)

type GenerateSeedOpts struct {
	Seed    string
	Version randomizer.Version
}

func ParseGenerateSeedOpts(args []string) (GenerateSeedOpts, error) {
	opts := GenerateSeedOpts{}
	flags := flag.NewFlagSet("generate", flag.ExitOnError)
	flags.StringVar(&opts.Seed, "s", "", "random number generator seed passed to the randomizer")
	flags.Var(&opts.Version, "r", "randomizer version to use")
	flags.Parse(args)
	if opts.Seed == "" {
		if seed, err := seed.GenerateSeedString(32); err != nil {
			return opts, err
		} else {
			opts.Seed = seed
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
	// Check if the actual randomizer even exists.
	if _, err := svc.GetRandomizer(opts.Version); err != nil {
		fmt.Printf("Could not use randomizer %s: %s.\n", opts.Version.String(), err.Error())
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
	fmt.Printf("Generating seed %s with randomizer %s", opts.Seed, opts.Version.String())
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
