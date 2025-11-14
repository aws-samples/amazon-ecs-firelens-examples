package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkGenAndWriteLogEntry(b *testing.B) {
	// Setup - create dependencies needed by genAndWriteLogEntry
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Use io.Discard to avoid actual output during benchmarking
	encoder := json.NewEncoder(io.Discard)

	// Create payload buffer with typical size (7KB - overhead = ~6KB)
	payloadSize := max((7*1024)-overhead, minPayload)
	payloadBuf := make([]byte, payloadSize)

	// Reset timer to exclude setup time from benchmark
	b.ResetTimer()

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		genAndWriteLogEntry(encoder, rng, payloadBuf, 7, 0)
	}
}

// Benchmark with different payload sizes
func BenchmarkGenAndWriteLogEntry_SmallPayload(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	encoder := json.NewEncoder(io.Discard)

	// Small payload: 1KB
	payloadSize := max((1*1024)-overhead, minPayload)
	payloadBuf := make([]byte, payloadSize)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		genAndWriteLogEntry(encoder, rng, payloadBuf, 1, 0)
	}
}

func BenchmarkGenAndWriteLogEntry_LargePayload(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	encoder := json.NewEncoder(io.Discard)

	// Large payload: 20KB
	payloadSize := max((20*1024)-overhead, minPayload)
	payloadBuf := make([]byte, payloadSize)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		genAndWriteLogEntry(encoder, rng, payloadBuf, 20, 0)
	}
}

// Memory allocation benchmark
func BenchmarkGenAndWriteLogEntry_Allocs(b *testing.B) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	encoder := json.NewEncoder(io.Discard)
	payloadSize := max((7*1024)-overhead, minPayload)
	payloadBuf := make([]byte, payloadSize)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		genAndWriteLogEntry(encoder, rng, payloadBuf, 7, 0)
	}
}
