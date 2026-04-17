package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

var testCols = []output.ColumnDef{
	{Header: "ID", Path: "id"},
	{Header: "NAME", Path: "name"},
}

type sampleItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func renderTo(t *testing.T, v any, opts output.RenderOpts) string {
	t.Helper()
	var buf bytes.Buffer
	if err := output.Render(&buf, v, opts); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	return buf.String()
}

// ─── Table ──────────────────────────────────────────────────────────────────

func TestRender_Table_SingleObject(t *testing.T) {
	item := sampleItem{ID: "1", Name: "Alice"}
	out := renderTo(t, item, output.RenderOpts{Format: output.FormatTable, Cols: testCols})
	if !strings.Contains(out, "ID") || !strings.Contains(out, "NAME") {
		t.Errorf("missing headers in: %q", out)
	}
	if !strings.Contains(out, "1") || !strings.Contains(out, "Alice") {
		t.Errorf("missing data in: %q", out)
	}
}

func TestRender_Table_Slice(t *testing.T) {
	items := []sampleItem{{ID: "1", Name: "Alice"}, {ID: "2", Name: "Bob"}}
	out := renderTo(t, items, output.RenderOpts{Format: output.FormatTable, Cols: testCols})
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "Bob") {
		t.Errorf("missing rows in: %q", out)
	}
	// Header should appear exactly once.
	if count := strings.Count(out, "ID"); count != 1 {
		t.Errorf("expected header once, got %d times", count)
	}
}

// ─── JSON ────────────────────────────────────────────────────────────────────

func TestRender_JSON_Pretty(t *testing.T) {
	item := sampleItem{ID: "x", Name: "Y"}
	out := renderTo(t, item, output.RenderOpts{Format: output.FormatJSON, Pretty: true})
	if !strings.Contains(out, "\n") {
		t.Error("expected indented JSON")
	}
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got["id"] != "x" {
		t.Errorf("id = %v", got["id"])
	}
}

func TestRender_JSON_Compact(t *testing.T) {
	item := sampleItem{ID: "x", Name: "Y"}
	out := renderTo(t, item, output.RenderOpts{Format: output.FormatJSON, Pretty: false})
	// Compact JSON should be on one line (plus newline from encoder).
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected compact JSON on 1 line, got %d lines", len(lines))
	}
}

// ─── NDJSON ──────────────────────────────────────────────────────────────────

func TestRender_NDJSON_Slice(t *testing.T) {
	items := []sampleItem{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}
	out := renderTo(t, items, output.RenderOpts{Format: output.FormatNDJSON})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 NDJSON lines, got %d: %q", len(lines), out)
	}
	for i, line := range lines {
		var obj map[string]any
		if err := json.Unmarshal([]byte(line), &obj); err != nil {
			t.Errorf("line %d not valid JSON: %v", i, err)
		}
	}
}

func TestRender_NDJSON_SingleObject(t *testing.T) {
	item := sampleItem{ID: "1", Name: "A"}
	out := renderTo(t, item, output.RenderOpts{Format: output.FormatNDJSON})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 NDJSON line for single object, got %d", len(lines))
	}
}

// ─── YAML ────────────────────────────────────────────────────────────────────

func TestRender_YAML(t *testing.T) {
	item := sampleItem{ID: "abc", Name: "Foo"}
	out := renderTo(t, item, output.RenderOpts{Format: output.FormatYAML})
	if !strings.Contains(out, "id: abc") {
		t.Errorf("expected YAML key 'id: abc' in %q", out)
	}
	if !strings.Contains(out, "name: Foo") {
		t.Errorf("expected YAML key 'name: Foo' in %q", out)
	}
}

// ─── TSV ─────────────────────────────────────────────────────────────────────

func TestRender_TSV(t *testing.T) {
	items := []sampleItem{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}
	out := renderTo(t, items, output.RenderOpts{Format: output.FormatTSV, Cols: testCols})
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 { // header + 2 data rows
		t.Fatalf("expected 3 TSV lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "\t") {
		t.Errorf("header missing tab: %q", lines[0])
	}
}

// ─── Plaintext ───────────────────────────────────────────────────────────────

func TestRender_Plaintext_String(t *testing.T) {
	out := renderTo(t, "hello world", output.RenderOpts{Format: output.FormatPlaintext})
	if strings.TrimSpace(out) != "hello world" {
		t.Errorf("got %q", out)
	}
	// No ANSI codes.
	if strings.Contains(out, "\x1b[") {
		t.Error("unexpected ANSI escape in plaintext")
	}
}

func TestRender_Plaintext_StringSlice(t *testing.T) {
	// Pass as []string — toAny marshals to []interface{} with string elements.
	out := renderTo(t, []string{"line1", "line2"}, output.RenderOpts{Format: output.FormatPlaintext})
	if !strings.Contains(out, "line1") || !strings.Contains(out, "line2") {
		t.Errorf("missing lines in %q", out)
	}
}

func TestRender_Plaintext_Struct(t *testing.T) {
	item := sampleItem{ID: "x", Name: "Y"}
	out := renderTo(t, item, output.RenderOpts{Format: output.FormatPlaintext})
	// Struct falls through to JSON marshal.
	if !strings.Contains(out, "x") {
		t.Errorf("expected struct content in %q", out)
	}
}

// ─── JQ filter ───────────────────────────────────────────────────────────────

func TestRender_JQ(t *testing.T) {
	item := sampleItem{ID: "42", Name: "Test"}
	out := renderTo(t, item, output.RenderOpts{
		Format: output.FormatJSON,
		JQ:     ".name",
		Pretty: false,
	})
	// .name should select the name string.
	trimmed := strings.TrimSpace(out)
	if trimmed != `"Test"` {
		t.Errorf("jq .name = %q, want \"Test\"", trimmed)
	}
}

func TestRender_JQ_Slice(t *testing.T) {
	items := []sampleItem{{ID: "1", Name: "A"}, {ID: "2", Name: "B"}}
	out := renderTo(t, items, output.RenderOpts{
		Format: output.FormatJSON,
		JQ:     ".[].id",
	})
	if !strings.Contains(out, `"1"`) || !strings.Contains(out, `"2"`) {
		t.Errorf("expected both ids in %q", out)
	}
}

func TestRender_JQ_InvalidExpr(t *testing.T) {
	var buf bytes.Buffer
	err := output.Render(&buf, sampleItem{}, output.RenderOpts{
		Format: output.FormatJSON,
		JQ:     "<<<invalid>>>",
	})
	if err == nil {
		t.Error("expected error for invalid jq expression")
	}
}

// ─── --fields projection ─────────────────────────────────────────────────────

func TestRender_Fields_SingleField(t *testing.T) {
	item := sampleItem{ID: "abc", Name: "keep"}
	out := renderTo(t, item, output.RenderOpts{
		Format: output.FormatJSON,
		Fields: "id",
	})
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := got["id"]; !ok {
		t.Error("expected 'id' field")
	}
	if _, ok := got["name"]; ok {
		t.Error("expected 'name' to be projected out")
	}
}

func TestRender_Fields_MultipleFields(t *testing.T) {
	item := map[string]any{"id": "1", "name": "Alice", "extra": "drop"}
	out := renderTo(t, item, output.RenderOpts{
		Format: output.FormatJSON,
		Fields: "id,name",
	})
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := got["extra"]; ok {
		t.Error("expected 'extra' to be projected out")
	}
	if got["id"] != "1" || got["name"] != "Alice" {
		t.Errorf("unexpected content: %v", got)
	}
}

func TestRender_Fields_DottedPath(t *testing.T) {
	nested := map[string]any{
		"user": map[string]any{
			"email": "a@b.com",
			"role":  "admin",
		},
		"other": "drop",
	}
	out := renderTo(t, nested, output.RenderOpts{
		Format: output.FormatJSON,
		Fields: "user.email",
	})
	var got map[string]any
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if _, ok := got["other"]; ok {
		t.Error("expected 'other' to be projected out")
	}
	userMap, ok := got["user"].(map[string]any)
	if !ok {
		t.Fatalf("user not a map: %v", got["user"])
	}
	if userMap["email"] != "a@b.com" {
		t.Errorf("email = %v", userMap["email"])
	}
	if _, ok := userMap["role"]; ok {
		t.Error("expected 'role' to be projected out of user")
	}
}
