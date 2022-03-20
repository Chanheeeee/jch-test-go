package model

import (
	"encoding/json"
)

type JSONPatchEntry struct {
	OP    string          `json:"op"`
	Path  string          `json:"path"`
	Value json.RawMessage `json:"value,omitempty"`
}