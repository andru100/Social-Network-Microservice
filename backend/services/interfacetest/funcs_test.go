package test

import (
    "testing"
)

func BenchmarkTes(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Test(10)
    }
}

// run command is go test -run ^$ -bench .