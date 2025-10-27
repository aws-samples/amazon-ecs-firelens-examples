package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	minPayload = 100
	overhead   = 225
)

var (
	// Pre-computed values
	charsetLen = len(charset)
	encoderMu  sync.Mutex

	// Static values that don't change
	staticFields = map[string]interface{}{
		"service":       "analytics",
		"event_type":    "video_view",
		"quality_level": "1080p",
		"device_type":   "smart_tv",
	}
)

type Config struct {
	Rate                 float64
	SizeKB               int
	BurstSizeKB          int
	BurstIntervalSeconds int
	Output               string
	RotateMaxSize        int
	RotateMaxBackups     int
	RotateMaxAge         int
}

func main() {
	// Parse environment variables
	config, err := parseConfig()
	if err != nil {
		panic(err)
	}

	payloadSize := max((config.SizeKB*1024)-overhead, minPayload)
	burstPayloadSize := max((config.BurstSizeKB*1024)-overhead, minPayload)

	burstInterval := time.Second * time.Duration(config.BurstIntervalSeconds)
	fmt.Printf("Starting SajaMediaGroup Analytics Log Generator...\n")
	fmt.Printf("Rate: %.1f logs/second\n", config.Rate)
	fmt.Printf("Size: ~%dKB per log\n", config.SizeKB)
	fmt.Printf("Burst Size: ~%dKB per log (burst every %s)\n", config.BurstSizeKB, burstInterval.String())
	fmt.Printf("Output File: %s\n", config.Output)

	interval := time.Duration(float64(time.Second) / config.Rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	burstTicker := time.NewTicker(burstInterval)
	defer burstTicker.Stop()
	var encoder *json.Encoder

	if config.Output != "STDOUT" {
		// Use lumberjack for log rotation
		logger := &lumberjack.Logger{
			Filename:   config.Output,
			MaxSize:    config.RotateMaxSize,    // MB
			MaxBackups: config.RotateMaxBackups, // Keep backup files
			MaxAge:     config.RotateMaxAge,     // Days (0 = don't delete based on age)
			Compress:   false,                   // Don't compress rotated files
		}
		encoder = json.NewEncoder(logger)
	} else {
		encoder = json.NewEncoder(os.Stdout)
	}
	for i := 0; i < 4; i++ {
		w := &worker{
			encoder:     encoder,
			rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
			buf:         make([]byte, payloadSize),
			burstBuf:    make([]byte, burstPayloadSize),
			ticker:      ticker,
			burstTicker: burstTicker,
		}
		go w.run()
	}
	w := &worker{
		encoder:     encoder,
		rng:         rand.New(rand.NewSource(time.Now().UnixNano())),
		buf:         make([]byte, payloadSize),
		burstBuf:    make([]byte, burstPayloadSize),
		ticker:      ticker,
		burstTicker: burstTicker,
	}
	w.run()
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func parseConfig() (*Config, error) {
	rate, err := strconv.ParseFloat(getEnv("LOG_RATE_PER_SECOND", "10"), 64)
	if err != nil {
		return nil, err
	}

	sizeKB, err := strconv.Atoi(getEnv("LOG_SIZE_KB", "7"))
	if err != nil {
		return nil, err
	}

	burstSizeKB := sizeKB
	if v := os.Getenv("LOG_SIZE_BURST_KB"); v != "" {
		burstSizeKB, err = strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
	}

	burstIntervalSeconds, err := strconv.Atoi(getEnv("BURST_INTERVAL_SECONDS", "15"))
	if err != nil {
		return nil, err
	}

	rotateMaxSize, err := strconv.Atoi(getEnv("LOG_ROTATE_MAX_SIZE", "100"))
	if err != nil {
		return nil, err
	}

	rotateMaxBackups, err := strconv.Atoi(getEnv("LOG_ROTATE_MAX_BACKUPS", "5"))
	if err != nil {
		return nil, err
	}

	rotateMaxAge, err := strconv.Atoi(getEnv("LOG_ROTATE_MAX_AGE", "1"))
	if err != nil {
		return nil, err
	}

	output := getEnv("OUTPUT_FILE", "STDOUT")
	return &Config{
		Rate:                 rate,
		SizeKB:               sizeKB,
		BurstSizeKB:          burstSizeKB,
		BurstIntervalSeconds: burstIntervalSeconds,
		Output:               output,
		RotateMaxSize:        rotateMaxSize,
		RotateMaxBackups:     rotateMaxBackups,
		RotateMaxAge:         rotateMaxAge,
	}, nil
}

func randomString(rng *rand.Rand, buf []byte) []byte {
	for i := 0; i < cap(buf); i++ {
		buf[i] = charset[rng.Intn(charsetLen)]
	}
	return buf
}

type worker struct {
	encoder     *json.Encoder
	rng         *rand.Rand
	buf         []byte
	burstBuf    []byte
	ticker      *time.Ticker
	burstTicker *time.Ticker
}

func (w *worker) run() {
	for range w.ticker.C {
		select {
		case <-w.burstTicker.C:
			genAndWriteLogEntry(w.encoder, w.rng, w.burstBuf)
		default:
			genAndWriteLogEntry(w.encoder, w.rng, w.buf)
		}
	}
}

func genAndWriteLogEntry(encoder *json.Encoder, rng *rand.Rand, payloadBuf []byte) {
	// Build log entry
	log := map[string]interface{}{
		"service":       staticFields["service"],
		"user_id":       strconv.Itoa(rng.Intn(90000) + 10000),
		"event_type":    staticFields["event_type"],
		"video_id":      strconv.Itoa(rng.Intn(90000) + 100000),
		"quality_level": staticFields["quality_level"],
		"device_type":   staticFields["device_type"],
		"metadata": map[string]interface{}{
			"buffer_events": rng.Intn(11),
			"seek_events":   rng.Intn(6),
			"cdn_edge":      "sea-01",
		},
		"large_payload": string(randomString(rng, payloadBuf)),
	}

	// Encode directly to stdout
	encoderMu.Lock()
	defer encoderMu.Unlock()
	_ = encoder.Encode(log)
}
