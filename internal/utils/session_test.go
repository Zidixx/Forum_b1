package utils

import "testing"

func TestGenerateUUID(t *testing.T) {
	id1 := GenerateUUID()
	id2 := GenerateUUID()

	if id1 == "" {
		t.Error("UUID should not be empty")
	}
	if len(id1) != 36 {
		t.Errorf("UUID should be 36 chars, got %d", len(id1))
	}
	if id1 == id2 {
		t.Error("two UUIDs should be different")
	}
}
