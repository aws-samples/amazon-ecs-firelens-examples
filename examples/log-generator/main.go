package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fluent/fluent-logger-golang/fluent"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	minPayload = 100
	overhead   = 275
)

type User struct {
	ID    string
	Email string
}

var (
	// Pre-computed values
	charsetLen = len(charset)
	writerMu   sync.Mutex

	// Buffer pool for JSON serialization to reduce allocations
	bufferPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// Sample users for random selection
	users = []User{
		{ID: "10001", Email: "alejandro_rosalez@amazondomains.com"},
		{ID: "10002", Email: "akua_mansa@amazondomains.com"},
		{ID: "10003", Email: "anacarolina_silva@amazondomains.com"},
		{ID: "10004", Email: "arnav_desai@amazondomains.com"},
		{ID: "10005", Email: "carlos_salazar@amazondomains.com"},
		{ID: "10006", Email: "diego_ramirez@amazondomains.com"},
		{ID: "10007", Email: "efua_owusu@amazondomains.com"},
		{ID: "10008", Email: "gil-dong_hong@amazondomains.com"},
		{ID: "10009", Email: "jane_doe@amazondomains.com"},
		{ID: "10010", Email: "ji-hoon_namgoong@amazondomains.com"},
		{ID: "10011", Email: "john_doe@amazondomains.com"},
		{ID: "10012", Email: "john_stiles@amazondomains.com"},
		{ID: "10013", Email: "jorge_souza@amazondomains.com"},
		{ID: "10014", Email: "kwaku_mensah@amazondomains.com"},
		{ID: "10015", Email: "kwesi_manu@amazondomains.com"},
		{ID: "10016", Email: "juan_li@amazondomains.com"},
		{ID: "10017", Email: "jie_liu@amazondomains.com"},
		{ID: "10018", Email: "marcia_oliveria@amazondomains.com"},
		{ID: "10019", Email: "maria_garcia@amazondomains.com"},
		{ID: "10020", Email: "martha_rivera@amazondomains.com"},
		{ID: "10021", Email: "mary_major@amazondomains.com"},
		{ID: "10022", Email: "mateo_jackson@amazondomains.com"},
		{ID: "10023", Email: "nikhil_jayashankar@amazondomains.com"},
		{ID: "10024", Email: "nikki_wolf@amazondomains.com"},
		{ID: "10025", Email: "pat_candella@amazondomains.com"},
		{ID: "10026", Email: "paulo_santos@amazondomains.com"},
		{ID: "10027", Email: "richard_roe@amazondomains.com"},
		{ID: "10028", Email: "saanvi_sarkar@amazondomains.com"},
		{ID: "10029", Email: "shirley_rodriguez@amazondomains.com"},
		{ID: "10030", Email: "sofia_martinez@amazondomains.com"},
		{ID: "10031", Email: "soo-jin_ki@amazondomains.com"},
		{ID: "10032", Email: "terry_whitlock@amazondomains.com"},
		{ID: "10033", Email: "xiulan_wang@amazondomains.com"},
		{ID: "10034", Email: "wei_zhang@amazondomains.com"},
	}
)

type Config struct {
	Rate                 float64
	SizeKB               int
	ExtraSizeKB          int
	BurstSizeKB          int
	BurstIntervalSeconds int
	Output               string
	RotateMaxSize        int
	RotateMaxBackups     int
	RotateMaxAge         int
	FluentHost           string
	FluentPort           int
	FluentTag            string
	UdpHost              string
	UdpPort              int
}

func main() {
	// Parse environment variables
	config, err := parseConfig()
	if err != nil {
		panic(err)
	}

	// Calculate max possible payload size for buffer allocation
	maxSizeKB := config.SizeKB + config.ExtraSizeKB
	payloadSize := max((maxSizeKB*1024)-overhead, minPayload)
	burstPayloadSize := max((config.BurstSizeKB*1024)-overhead, minPayload)

	burstInterval := time.Second * time.Duration(config.BurstIntervalSeconds)
	fmt.Printf("Starting AnyCompany Media Group Analytics Log Generator...\n")
	fmt.Printf("Rate: %.1f logs/second\n", config.Rate)
	if config.ExtraSizeKB > 0 {
		fmt.Printf("Size: %d-%dKB per log (variable)\n", config.SizeKB, config.SizeKB+config.ExtraSizeKB)
	} else {
		fmt.Printf("Size: ~%dKB per log\n", config.SizeKB)
	}
	fmt.Printf("Burst Size: ~%dKB per log (burst every %s)\n", config.BurstSizeKB, burstInterval.String())
	fmt.Printf("Output File: %s\n", config.Output)

	interval := time.Duration(float64(time.Second) / config.Rate)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	burstTicker := time.NewTicker(burstInterval)
	defer burstTicker.Stop()
	var writer io.Writer

	// fluentLogger only gets created for the FORWARD output
	var fluentLogger *fluent.Fluent
	switch config.Output {
	case "FORWARD":
		// Use Fluent-bit forwarder
		writer = io.Discard
		fluentLogger = newFluentClient(config)
		defer fluentLogger.Close()
		fmt.Printf("Fluent logger configured: %s:%d (tag: %s)\n", config.FluentHost, config.FluentPort, config.FluentTag)
	case "UDP":
		// Use UDP output
		udpAddr := fmt.Sprintf("%s:%d", config.UdpHost, config.UdpPort)
		conn, err := net.Dial("udp", udpAddr)
		if err != nil {
			panic(fmt.Sprintf("Failed to create UDP connection to %s: %v", udpAddr, err))
		}
		defer conn.Close()
		writer = conn
		fmt.Printf("UDP output configured: %s\n", udpAddr)
	case "STDOUT":
		writer = os.Stdout
	default:
		// default assumes a file-path was provided for log-file output.
		// Use lumberjack for log rotation
		writer = &lumberjack.Logger{
			Filename:   config.Output,
			MaxSize:    config.RotateMaxSize,    // MB
			MaxBackups: config.RotateMaxBackups, // Keep backup files
			MaxAge:     config.RotateMaxAge,     // Days (0 = don't delete based on age)
			Compress:   false,                   // Don't compress rotated files
		}
	}
	// create 5 workers to run async
	for i := 0; i < 5; i++ {
		w := &worker{
			writer:       writer,
			fluentLogger: fluentLogger,
			rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
			buf:          make([]byte, payloadSize),
			burstBuf:     make([]byte, burstPayloadSize),
			ticker:       ticker,
			burstTicker:  burstTicker,
			config:       config,
		}
		go w.run()
	}
	// 1 worker runs synchronously
	w := &worker{
		writer:       writer,
		fluentLogger: fluentLogger,
		rng:          rand.New(rand.NewSource(time.Now().UnixNano())),
		buf:          make([]byte, payloadSize),
		burstBuf:     make([]byte, burstPayloadSize),
		ticker:       ticker,
		burstTicker:  burstTicker,
		config:       config,
	}
	w.run()
}

type worker struct {
	writer       io.Writer
	fluentLogger *fluent.Fluent
	rng          *rand.Rand
	buf          []byte
	burstBuf     []byte
	ticker       *time.Ticker
	burstTicker  *time.Ticker
	config       *Config
}

func (w *worker) run() {
	for range w.ticker.C {
		select {
		case <-w.burstTicker.C:
			genAndWriteLogEntry(w.writer, w.fluentLogger, w.config, w.rng, w.burstBuf, w.config.BurstSizeKB, 0)
		default:
			genAndWriteLogEntry(w.writer, w.fluentLogger, w.config, w.rng, w.buf, w.config.SizeKB, w.config.ExtraSizeKB)
		}
	}
}

func genAndWriteLogEntry(
	writer io.Writer,
	fluentLogger *fluent.Fluent,
	config *Config,
	rng *rand.Rand,
	payloadBuf []byte,
	baseSizeKB,
	extraSizeKB int,
) {
	log := buildLogEntry(rng, payloadBuf, baseSizeKB, extraSizeKB)

	if fluentLogger != nil {
		// Use Fluent logger to send log map directly
		if err := fluentLogger.Post(config.FluentTag, log); err != nil {
			// Log error but continue
			fmt.Fprintf(os.Stderr, "Failed to post to fluent: %v\n", err)
		}
		return
	}

	// Get buffer from pool and reset it
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(log); err != nil {
		return // Skip this log entry on encode error
	}

	// Only acquire lock for the write operation
	writerMu.Lock()
	defer writerMu.Unlock()
	_, _ = writer.Write(buf.Bytes())
}

func buildLogEntry(rng *rand.Rand, payloadBuf []byte, baseSizeKB, extraSizeKB int) map[string]any {
	// Calculate actual payload size with extra random size
	actualSizeKB := baseSizeKB
	if extraSizeKB > 0 {
		// Apply random extra size: baseSizeKB + (0 to extraSizeKB)
		extraSize := rng.Intn(extraSizeKB + 1) // Range: 0 to extraSizeKB (inclusive)
		actualSizeKB = baseSizeKB + extraSize
	}

	actualPayloadSize := max((actualSizeKB*1024)-overhead, minPayload)
	// Use a slice of the buffer up to the actual size needed
	actualBuf := payloadBuf[:min(actualPayloadSize, len(payloadBuf))]

	// Select random user
	user := users[rng.Intn(len(users))]

	// Determine event type (90% ongoing, 10% final)
	var eventType string
	var durationS int
	if rng.Intn(10) < 9 {
		// 90% - video_view_ongoing: 1s to 3h (1 to 10,800 seconds)
		eventType = "video_view_ongoing"
		durationS = rng.Intn(10800) + 1
	} else {
		// 10% - video_view_final: 20m to 3h (1,200 to 10,800 seconds)
		eventType = "video_view_final"
		durationS = rng.Intn(9601) + 1200
	}

	// Build log entry
	return map[string]any{
		"service":       "analytics",
		"user_id":       user.ID,
		"user_email":    user.Email,
		"event_type":    eventType,
		"video_id":      strconv.Itoa(rng.Intn(90000) + 100000),
		"duration_s":    durationS,
		"quality_level": "1080p",
		"device_type":   "smart_tv",
		"metadata": map[string]any{
			"buffer_events": rng.Intn(11),
			"seek_events":   rng.Intn(6),
			"cdn_edge":      "sea-01",
		},
		"large_payload": string(randomString(rng, actualBuf)),
	}
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

	extraSizeKB, err := strconv.Atoi(getEnv("LOG_SIZE_EXTRA_KB", "0"))
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

	rotateMaxSize, err := strconv.Atoi(getEnv("LOG_ROTATE_MAX_SIZE_MB", "100"))
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

	fluentHost := getEnv("FLUENT_HOST", "localhost")

	fluentPort, err := strconv.Atoi(getEnv("FLUENT_PORT", "24224"))
	if err != nil {
		return nil, err
	}

	fluentTag := getEnv("FLUENT_TAG", "video.analytics")

	udpHost := getEnv("UDP_HOST", "localhost")

	udpPort, err := strconv.Atoi(getEnv("UDP_PORT", "5170"))
	if err != nil {
		return nil, err
	}

	return &Config{
		Rate:                 rate,
		SizeKB:               sizeKB,
		ExtraSizeKB:          extraSizeKB,
		BurstSizeKB:          burstSizeKB,
		BurstIntervalSeconds: burstIntervalSeconds,
		Output:               output,
		RotateMaxSize:        rotateMaxSize,
		RotateMaxBackups:     rotateMaxBackups,
		RotateMaxAge:         rotateMaxAge,
		FluentHost:           fluentHost,
		FluentPort:           fluentPort,
		FluentTag:            fluentTag,
		UdpHost:              udpHost,
		UdpPort:              udpPort,
	}, nil
}

func randomString(rng *rand.Rand, buf []byte) []byte {
	for i := 0; i < len(buf); i++ {
		buf[i] = charset[rng.Intn(charsetLen)]
	}
	return buf
}

func newFluentClient(config *Config) *fluent.Fluent {
	if config.Output != "FORWARD" {
		return nil
	}

	fluentConfig := fluent.Config{
		FluentHost: config.FluentHost,
		FluentPort: config.FluentPort,
	}
	client, err := fluent.New(fluentConfig)
	if err != nil {
		panic(fmt.Sprintf("Failed to create fluent logger: %v", err))
	}
	return client
}
