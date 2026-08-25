package zerolog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"testing"
)

type syncBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}

func (s *syncBuffer) Bytes() []byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]byte, s.buf.Len())
	copy(cp, s.buf.Bytes())
	return cp
}

func TestConcurrentContextDerivation(t *testing.T) {
	var buf syncBuffer
	base := New(&buf).With().Str("service", "api").Logger()

	const n = 500
	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			reqID := fmt.Sprintf("req-%d", id)
			sub := base.With().Str("request_id", reqID).Logger()
			sub.Info().Msg("handled")
		}(i)
	}

	wg.Wait()

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	if len(lines) != n {
		t.Fatalf("expected %d log lines, got %d", n, len(lines))
	}

	foundReqIDs := make(map[string]bool)
	for idx, line := range lines {
		var entry map[string]interface{}
		if err := json.Unmarshal(line, &entry); err != nil {
			t.Fatalf("line %d is not valid JSON: %s (err: %v)", idx, string(line), err)
		}
		if entry["service"] != "api" {
			t.Errorf("line %d: expected service=api, got %v", idx, entry["service"])
		}
		if entry["message"] != "handled" {
			t.Errorf("line %d: expected message=handled, got %v", idx, entry["message"])
		}
		reqID, ok := entry["request_id"].(string)
		if !ok || reqID == "" {
			t.Errorf("line %d: missing or invalid request_id: %v", idx, entry["request_id"])
		} else {
			foundReqIDs[reqID] = true
		}
	}

	if len(foundReqIDs) != n {
		t.Errorf("expected %d unique request_ids, got %d", n, len(foundReqIDs))
	}
}

func TestContextIsolation(t *testing.T) {
	var buf bytes.Buffer
	root := New(&buf).With().Str("root", "val").Logger()
	child1 := root.With().Str("child", "1").Logger()
	child2 := root.With().Str("child", "2").Logger()

	buf.Reset()
	child1.Info().Msg("c1")
	var entry1 map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &entry1)
	if entry1["root"] != "val" || entry1["child"] != "1" {
		t.Errorf("unexpected child1 entry: %v", entry1)
	}

	buf.Reset()
	child2.Info().Msg("c2")
	var entry2 map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &entry2)
	if entry2["root"] != "val" || entry2["child"] != "2" {
		t.Errorf("unexpected child2 entry: %v", entry2)
	}
}

func BenchmarkContextDerivation(b *testing.B) {
	logger := New(io.Discard).With().Str("service", "api").Logger()
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = logger.With().Str("request_id", "12345").Logger()
	}
}
