package core

import (
	"encoding/json"
	"testing"
)

func TestExtractPanelFallbacksToBytesNormalizesFallbacks(t *testing.T) {
	raw := json.RawMessage(`[
		{"name":"h2","alpn":"h2","path":"/h2","type":"tcp","dest":8080,"xver":"1"},
		{"alpn":"http/1.1","dest":"127.0.0.1:80","xver":2},
		{"alpn":"missing-dest"}
	]`)

	out, err := extractPanelFallbacksToBytes(raw)
	if err != nil {
		t.Fatalf("extractPanelFallbacksToBytes() error = %v", err)
	}

	var got []map[string]interface{}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len(got) = %d, want 2", len(got))
	}
	if got[0]["dest"] != "8080" {
		t.Fatalf("got[0][dest] = %v, want 8080", got[0]["dest"])
	}
	if got[0]["xver"] != float64(1) {
		t.Fatalf("got[0][xver] = %v, want 1", got[0]["xver"])
	}
	if got[1]["dest"] != "127.0.0.1:80" {
		t.Fatalf("got[1][dest] = %v, want 127.0.0.1:80", got[1]["dest"])
	}
}

func TestExtractPanelFallbacksToBytesRejectsInvalidXver(t *testing.T) {
	raw := json.RawMessage(`[{"dest":"127.0.0.1:80","xver":"bad"}]`)

	if _, err := extractPanelFallbacksToBytes(raw); err == nil {
		t.Fatal("extractPanelFallbacksToBytes() error = nil, want error")
	}
}

func TestExtractPanelFallbacksToBytesRejectsFractionalXver(t *testing.T) {
	raw := json.RawMessage(`[{"dest":"127.0.0.1:80","xver":1.5}]`)

	if _, err := extractPanelFallbacksToBytes(raw); err == nil {
		t.Fatal("extractPanelFallbacksToBytes() error = nil, want error")
	}
}
