package output_test

import (
	"testing"

	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

func TestProject_NoFields(t *testing.T) {
	v := map[string]any{"a": 1, "b": 2}
	got := output.Project(v, nil)
	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", got)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 keys, got %d", len(m))
	}
}

func TestProject_TopLevel(t *testing.T) {
	v := map[string]any{"a": 1, "b": 2, "c": 3}
	got := output.Project(v, []string{"a", "c"})
	m, ok := got.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", got)
	}
	if _, has := m["b"]; has {
		t.Error("expected 'b' to be projected out")
	}
	if m["a"] == nil || m["c"] == nil {
		t.Error("expected 'a' and 'c' present")
	}
}

func TestProject_DottedPath_TwoLevels(t *testing.T) {
	v := map[string]any{
		"user": map[string]any{
			"email": "x@y.com",
			"name":  "Alice",
		},
		"other": "drop",
	}
	got := output.Project(v, []string{"user.email"})
	m := got.(map[string]any)
	if _, has := m["other"]; has {
		t.Error("'other' should be projected out")
	}
	user := m["user"].(map[string]any)
	if user["email"] != "x@y.com" {
		t.Errorf("email = %v", user["email"])
	}
	if _, has := user["name"]; has {
		t.Error("'name' should be projected out")
	}
}

func TestProject_Slice(t *testing.T) {
	v := []any{
		map[string]any{"a": 1, "b": 2},
		map[string]any{"a": 3, "b": 4},
	}
	got := output.Project(v, []string{"a"})
	slice, ok := got.([]any)
	if !ok {
		t.Fatalf("expected slice, got %T", got)
	}
	for i, item := range slice {
		m := item.(map[string]any)
		if _, has := m["b"]; has {
			t.Errorf("item %d: 'b' should be projected out", i)
		}
	}
}

func TestProject_MissingField_Ignored(t *testing.T) {
	v := map[string]any{"a": 1}
	got := output.Project(v, []string{"a", "z"})
	m := got.(map[string]any)
	if _, has := m["z"]; has {
		t.Error("'z' should not appear since it's not in input")
	}
	if m["a"] == nil {
		t.Error("'a' should be present")
	}
}

func TestProject_NonMap_PassThrough(t *testing.T) {
	// Strings, numbers, booleans are returned unchanged.
	got := output.Project("hello", []string{"a"})
	if got != "hello" {
		t.Errorf("expected string passthrough, got %v", got)
	}
}

func TestProject_MixedSiblings(t *testing.T) {
	// Requesting "user" and "user.email" — "user" wins (include whole field).
	v := map[string]any{
		"user": map[string]any{
			"email": "a@b.com",
			"role":  "admin",
		},
	}
	got := output.Project(v, []string{"user", "user.email"})
	m := got.(map[string]any)
	user := m["user"].(map[string]any)
	// Both sub-keys should be present because "user" was requested whole.
	if user["role"] == nil {
		t.Error("expected 'role' present when 'user' requested as whole")
	}
}
