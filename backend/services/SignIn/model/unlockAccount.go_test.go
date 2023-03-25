package model

import (
    "testing"
)

func BenchmarkTes(b *testing.B) {
    a:= "andrew"
    for i := 0; i < b.N; i++ {
        UnlockAccount(&a)
    }
}

// run command is go test -run ^$ -bench .