package runner

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/HenryNg101/docker-monitoring-agent/internal/bootstrap"
	"github.com/HenryNg101/docker-monitoring-agent/internal/metrics"
	"github.com/HenryNg101/docker-monitoring-agent/internal/security"
)

// TODO: Review all these content later
const (
	metricsInterval  = 30 * time.Second
	rotationInterval = 24 * time.Hour
	flushInterval    = 180 * time.Second

	maxBatchSize  = 500
	maxBufferSize = 5000

	bufferFile = "/app/data/buffer.json"
)

type Runner struct {
	cfg    *bootstrap.Config
	sec    *bootstrap.Secret
	client *http.Client

	buffer []metrics.MetricMessage
	mu     sync.Mutex
}

func NewRunner(cfg *bootstrap.Config, sec *bootstrap.Secret, client *http.Client) *Runner {
	return &Runner{
		cfg:    cfg,
		sec:    sec,
		client: client,
	}
}

// add small randomness to avoid all agents syncing together
func jitter(base time.Duration) time.Duration {
	return base + time.Duration(rand.Int63n(int64(base/5)))
}

func (r *Runner) Start() {
	r.loadBuffer()

	collectTicker := time.NewTicker(jitter(metricsInterval))
	flushTicker := time.NewTicker(jitter(flushInterval))
	rotationTicker := time.NewTicker(jitter(rotationInterval))

	defer collectTicker.Stop()
	defer flushTicker.Stop()
	defer rotationTicker.Stop()
	defer r.persistBuffer()

	for {
		select {
		case <-collectTicker.C:
			r.collect()

		case <-flushTicker.C:
			r.flush()

		case <-rotationTicker.C:
			r.handleRotation()

		default:
			// soft limits without locking too often
			r.mu.Lock()
			if len(r.buffer) > maxBufferSize {
				// drop oldest
				r.buffer = r.buffer[len(r.buffer)-maxBufferSize:]
			}
			shouldFlush := len(r.buffer) >= maxBatchSize
			r.mu.Unlock()

			if shouldFlush {
				r.flush()
			}

			time.Sleep(300 * time.Millisecond)
		}
	}
}

func (r *Runner) collect() {
	messages, err := collectMetrics(r.cfg.ServerID)
	if err != nil {
		log.Println("[Agent] collect error:", err)
		return
	}

	r.mu.Lock()
	r.buffer = append(r.buffer, messages...)
	r.mu.Unlock()
}

func (r *Runner) flush() {
	// 1. Take snapshot safely
	r.mu.Lock()
	if len(r.buffer) == 0 {
		r.mu.Unlock()
		return
	}

	// copy buffer → snapshot
	data := make([]metrics.MetricMessage, len(r.buffer))
	copy(data, r.buffer)

	// clear original buffer
	r.buffer = nil
	r.mu.Unlock()

	// 2. Marshal outside lock (slow op)
	payload, err := json.Marshal(data)
	if err != nil {
		log.Println("[Agent] marshal error:", err)
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		r.cfg.APIURL+"/agent/metrics",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		log.Println("[Agent] request error:", err)
		return
	}

	// 3. decrypt API key
	key := security.DeriveKey(r.cfg.InstanceID)
	rawKey, err := security.Decrypt(r.sec.APIKeyEncrypted, key)
	if err != nil {
		log.Println("[Agent] decrypt failed:", err)
		return
	}

	req.Header.Set("X-Agent-API-Key", rawKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		log.Println("[Agent] send failed:", err)

		// ❗ restore buffer on failure
		r.restoreBuffer(data)
		r.persistBuffer()
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		log.Println("[Agent] auth failed → rotating key")

		_ = bootstrap.RotateKey(r.cfg, r.sec)

		// ❗ restore buffer so we don't lose data
		r.restoreBuffer(data)
		return
	}

	log.Println("[Agent] metrics sent:", len(data))
}

func (r *Runner) restoreBuffer(data []metrics.MetricMessage) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// prepend back failed batch
	r.buffer = append(data, r.buffer...)
}

func (r *Runner) persistBuffer() {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := json.Marshal(r.buffer)
	if err != nil {
		return
	}

	_ = os.WriteFile(bufferFile, data, 0600)
}

func (r *Runner) loadBuffer() {
	data, err := os.ReadFile(bufferFile)
	if err != nil {
		return
	}

	var buf []metrics.MetricMessage
	if err := json.Unmarshal(data, &buf); err != nil {
		return
	}

	r.buffer = buf
}
