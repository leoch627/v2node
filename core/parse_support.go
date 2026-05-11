package core

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type panelFallback struct {
	Name string          `json:"name,omitempty"`
	Alpn string          `json:"alpn,omitempty"`
	Path string          `json:"path,omitempty"`
	Type string          `json:"type,omitempty"`
	Dest json.RawMessage `json:"dest,omitempty"`
	Xver json.RawMessage `json:"xver,omitempty"`
}

func extractPanelFallbacksToBytes(raw json.RawMessage) ([]byte, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var panelFBs []panelFallback
	if err := json.Unmarshal(raw, &panelFBs); err != nil {
		return nil, fmt.Errorf("unmarshal panel fallbacks: %w", err)
	}

	var parsedFBs []map[string]interface{}
	for i, fb := range panelFBs {
		m := make(map[string]interface{})
		if fb.Name != "" {
			m["name"] = fb.Name
		}
		if fb.Alpn != "" {
			m["alpn"] = fb.Alpn
		}
		if fb.Path != "" {
			m["path"] = fb.Path
		}
		if fb.Type != "" {
			m["type"] = fb.Type
		}

		dest, err := parsePanelFallbackDest(fb.Dest)
		if err != nil {
			return nil, fmt.Errorf("fallback[%d] dest: %w", i, err)
		}
		if dest != "" {
			m["dest"] = dest
		}
		if dest == "" {
			continue
		}

		xver, err := parsePanelFallbackXver(fb.Xver)
		if err != nil {
			return nil, fmt.Errorf("fallback[%d] xver: %w", i, err)
		}
		if xver > 0 {
			m["xver"] = xver
		}

		parsedFBs = append(parsedFBs, m)
	}
	out, err := json.Marshal(parsedFBs)
	if err != nil {
		return nil, fmt.Errorf("marshal panel fallbacks: %w", err)
	}
	return out, nil
}

func parsePanelFallbackDest(raw json.RawMessage) (string, error) {
	if len(raw) == 0 {
		return "", nil
	}
	var dstr string
	if err := json.Unmarshal(raw, &dstr); err == nil {
		return strings.TrimSpace(dstr), nil
	}

	var dint int
	if err := json.Unmarshal(raw, &dint); err == nil {
		return strconv.Itoa(dint), nil
	}
	return "", fmt.Errorf("must be a string or integer")
}

func parsePanelFallbackXver(raw json.RawMessage) (uint64, error) {
	if len(raw) == 0 {
		return 0, nil
	}

	var xv interface{}
	if err := json.Unmarshal(raw, &xv); err != nil {
		return 0, err
	}
	switch v := xv.(type) {
	case float64:
		if v < 0 || math.Trunc(v) != v {
			return 0, fmt.Errorf("must be a non-negative integer")
		}
		return uint64(v), nil
	case string:
		if v == "" {
			return 0, nil
		}
		xver, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0, err
		}
		return xver, nil
	default:
		return 0, fmt.Errorf("must be a string or integer")
	}
}
