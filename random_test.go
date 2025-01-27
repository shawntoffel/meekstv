package meekstv

import "testing"

func TestPseudoRandom(t *testing.T) {
	random := NewWichmannHillRandom(43, 10003, 20009)
	assertEqual(t, random.NextInt(3), 0)
	assertEqual(t, random.NextInt(3), 2)
	assertEqual(t, random.NextInt(3), 1)
}

func BenchmarkPseudoRandom(b *testing.B) {
	random := NewWichmannHillRandom(43, 10003, 20009)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		random.NextInt(3)
	}
}
