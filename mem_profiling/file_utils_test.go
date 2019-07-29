package mem_profiling

import (
	"testing"
)

func BenchmarkFindInFile(b *testing.B) {
	b.ResetTimer()
	count := 0
	for i := 0; i < b.N; i++ {
		count, _ = findInFile("Fermi paradox")
	}
	println("keyword:Fermi, count:", count)
}
