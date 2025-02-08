package python

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Interpreter interface {
	Version() (string, error)
}

// LocalInterpreter is a Python interpreter on the local filesystem.
type LocalInterpreter struct {
	path string
}

func pythonExecutableNames() []string {
	if runtime.GOOS == "windows" {
		return []string{"python.exe"}
	}
	return []string{"python3", "python"}
}

func findExecutable(path string, names ...string) (string, error) {
	for _, name := range names {
		filename := filepath.Join(path, name)
		if info, err := os.Stat(filename); err == nil && !info.IsDir() {
			return filename, nil
		}
	}
	return "", errors.New("not found")
}

func FindInterpreter() (*LocalInterpreter, error) {
	return FindInterpreterByName(pythonExecutableNames()...)
}

func FindInterpreterByName(names ...string) (*LocalInterpreter, error) {
	paths := strings.Split(os.Getenv("PATH"), string(os.PathListSeparator))
	for _, path := range paths {
		executable, err := findExecutable(path, names...)
		if err != nil {
			continue
		}
		return &LocalInterpreter{executable}, nil
	}
	return nil, errors.New("not found")
}

func (i *LocalInterpreter) Path() string {
	return i.path
}

func (i *LocalInterpreter) Command(arg ...string) *exec.Cmd {
	return exec.Command(i.Path(), arg...)
}

func (i *LocalInterpreter) Version() (string, error) {
	cmd := i.Command("--version")
	b, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	output := strings.TrimSpace(string(b))
	version := strings.TrimPrefix(output, "Python ")
	return version, nil
}
