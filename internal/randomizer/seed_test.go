package randomizer_test

import (
	"crypto/rand"
	"github.com/jessebrands/triforceblitz/internal/randomizer"
	"io"
	"testing"
)

// skipIfPrngUnavailable asserts that a PRNG is available on the test system.
func skipIfPrngUnavailable(t *testing.T) {
	b := make([]byte, 1)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		t.Skip("No PRNG available, skipping tests")
	}
}

func TestGenerateSeedString(t *testing.T) {
	skipIfPrngUnavailable(t)

	t.Run("must not error", func(t *testing.T) {
		if _, err := randomizer.GenerateSeedString(32); err != nil {
			t.Errorf("expected nil error, got %s", err.Error())
		}
	})

	t.Run("len must be equal to n", func(t *testing.T) {
		if a, _ := randomizer.GenerateSeedString(32); len(a) != 32 {
			t.Errorf("expected len to be 32, got %d", len(a))
		}
	})

	t.Run("output must not be equal", func(t *testing.T) {
		a, _ := randomizer.GenerateSeedString(32)
		b, _ := randomizer.GenerateSeedString(32)
		if a == b {
			t.Errorf("expected %s != %s", a, b)
		}
	})
}

func BenchmarkGenerateSeedString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randomizer.GenerateSeedString(32)
	}
}
