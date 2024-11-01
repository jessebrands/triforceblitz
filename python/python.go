package python

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Interpreter struct {
	path string
}

var defaultNames = []string{"python3", "python"}

// Find searches the system PATH for a Python interpreter.
//
// It seeks the system path for an executable in the following order:
//  1. python3
//  2. python
func Find() (*Interpreter, error) {
	return FindNames(defaultNames...)
}

// FindNames searches the system PATH for a Python interpreter.
func FindNames(files ...string) (*Interpreter, error) {
	for _, file := range files {
		interpreter, err := FindName(file)
		if err != nil {
			continue
		}
		return interpreter, nil
	}
	return nil, fmt.Errorf("could not find an interpreter")
}

// FindName searches the system PATH for a Python interpreter with a specific
// name.
func FindName(file string) (*Interpreter, error) {
	bin, err := exec.LookPath(file)
	if err != nil {
		return nil, err
	}
	path, err := filepath.Abs(bin)
	if err != nil {
		return nil, err
	}
	return &Interpreter{path}, nil
}

func (i *Interpreter) command(args ...string) *exec.Cmd {
	return exec.Command(i.path, args...)
}

func (i *Interpreter) Command(script string, args ...string) *exec.Cmd {
	args = append([]string{script}, args...)
	return i.command(args...)
}

func (i *Interpreter) Version() (string, error) {
	cmd := i.command("--version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	version := strings.TrimPrefix(string(out), "Python ")
	return strings.TrimSpace(version), nil
}

func (i *Interpreter) Path() string {
	return i.path 
}

