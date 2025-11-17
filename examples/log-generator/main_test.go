package main

import (
	"io"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkGenAndWriteLogEntry(b *testing.B) {
	// Setup - create dependencies needed by genAndWriteLogEntry
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	config := &Config{FluentTag: "test"}

	// Create payload buffer with typical size (7KB - overhead = ~6KB)
	payloadSize := max((7*1024)-overhead, minPayload)
	payloadBuf := make([]byte, payloadSize)

	// Reset timer to exclude setup time from benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		genAndWriteLogEntry(io.Discard, nil, config, rng, payloadBuf, 7, 0)
	}
}

// Benchmark with different payload sizes
func BenchmarkGenAndWriteLogEntry_SmallPayload(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	config := &Config{FluentTag: "test"}

	// Small payload: 1KB
	payloadSize := max((1*1024)-overhead, minPayload)
	payloadBuf := make([]byte, payloadSize)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		genAndWriteLogEntry(io.Discard, nil, config, rng, payloadBuf, 1, 0)
	}
}

func BenchmarkGenAndWriteLogEntry_LargePayload(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	config := &Config{FluentTag: "test"}

	// Large payload: 20KB
	payloadSize := max((20*1024)-overhead, minPayload)
	payloadBuf := make([]byte, payloadSize)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		genAndWriteLogEntry(io.Discard, nil, config, rng, payloadBuf, 20, 0)
	}
}

// Memory allocation benchmark
func BenchmarkGenAndWriteLogEntry_Allocs(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	config := &Config{FluentTag: "test"}
	payloadSize := max((7*1024)-overhead, minPayload)
	payloadBuf := make([]byte, payloadSize)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		genAndWriteLogEntry(io.Discard, nil, config, rng, payloadBuf, 7, 0)
	}
}
