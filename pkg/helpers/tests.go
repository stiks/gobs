package helpers

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"

	"github.com/google/uuid"
)

func UUIDFromString(t *testing.T, uStr string) uuid.UUID {
	id, err := uuid.Parse(uStr)
	if err != nil {
		if t == nil {
			log.Fatalf("Unable to get UUID from string: %s", err.Error())
		} else {
			t.Fatalf("Unable to get UUID from string: %s", err.Error())
		}
	}

	return id
}

func UUIDFromStringPointer(t *testing.T, uStr string) *uuid.UUID {
	id := UUIDFromString(t, uStr)

	return &id
}

func ObjectToByte(t *testing.T, obj interface{}) *bytes.Reader {
	b, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("Failed to marshal body: %v", err.Error())
	}

	return bytes.NewReader(b)
}
