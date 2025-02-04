package python_test

import (
	"errors"
	"github.com/jessebrands/triforceblitz/internal/python"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

var pythonVersionRegexp = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`)

// findPython searches for a Python interpreter on the system PATH. If no interpreter
// can be found, an empty string is returned. Check the error for more information.
func findPython() (string, error) {
	names := []string{"python.exe", "python3", "python"}
	for _, name := range names {
		path, err := exec.LookPath(name)
		if err == nil {
			log.Printf("Found Python executable at %s\n", path)
			return path, nil
		}
	}
	return "", errors.New("not found")
}

// findPythonOrSkip checks that a Python interpreter is installed and returns
// the path to the default system Python interpreter if found. This interpreter can
// be used to confirm test results later on.
//
// If it fails to find any interpreter, the test will be skipped.
func findPythonOrSkip(t *testing.T) string {
	t.Helper()
	path, err := findPython()
	if err != nil {
		t.Skipf("python not found in PATH")
	}
	return path
}

func TestFindInterpreter(t *testing.T) {
	findPythonOrSkip(t)
	interpreter, err := python.FindInterpreter()
	if err != nil {
		t.Fatal(err)
	}
	if interpreter == nil {
		t.Errorf("expected a non-nil interpreter")
	}
}

func TestFindInterpreterByNameFail(t *testing.T) {
	// In a long, long distant future, this might fail. :-)
	interpreter, err := python.FindInterpreterByName("this-should-always-fail")
	if err == nil {
		t.Errorf("expected a non-nil error")
	}
	if err != nil && !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected error message to contain 'not found'")
	}
	if interpreter != nil {
		t.Errorf("expected a nil interpreter")
	}
}

func TestLocalInterpreter_Path(t *testing.T) {
	path := findPythonOrSkip(t)
	interpreter, _ := python.FindInterpreter()
	if interpreter == nil {
		t.Errorf("expected a non-nil interpreter")
	}
	if interpreter.Path() != path {
		t.Errorf("expected interpreter path to be %s, got %s", path, interpreter.Path())
	}
}

func TestLocalInterpreter_Version(t *testing.T) {
	findPythonOrSkip(t)
	interpreter, _ := python.FindInterpreter()
	version, err := interpreter.Version()
	if err != nil {
		t.Errorf("expected a non-nil error")
	}
	if !pythonVersionRegexp.MatchString(version) {
		t.Errorf("expected version to match %s", pythonVersionRegexp.String())
	}
}
