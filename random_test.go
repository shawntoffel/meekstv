package meekstv

import "testing"

func TestPseudoRandom(t *testing.T) {
	random := NewWichmannHillRandom(43, 10003, 20009)
	result := random.NextInt(3)
	if result != 0 {
		t.Errorf("expected: %d, got: %d", 0, result)
	}
	result = random.NextInt(3)
	if result != 2 {
		t.Errorf("expected: %d, got: %d", 2, result)
	}
	result = random.NextInt(3)
	if result != 1 {
		t.Errorf("expected: %d, got: %d", 1, result)
	}
}

func BenchmarkPseudoRandom(b *testing.B) {
	random := NewWichmannHillRandom(43, 10003, 20009)
	for i := 0; i < b.N; i++ {
		random.NextInt(3)
	}
}
