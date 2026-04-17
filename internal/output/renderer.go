// Package output renders command responses across formats (table, JSON,
// NDJSON, YAML, TSV, plaintext) and applies post-processing (jq filter,
// client-side field projection).
package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/itchyny/gojq"
	"gopkg.in/yaml.v3"
)

// Format enumerates the output formats.
type Format string

const (
	FormatTable     Format = "table"
	FormatJSON      Format = "json"
	FormatNDJSON    Format = "ndjson"
	FormatYAML      Format = "yaml"
	FormatTSV       Format = "tsv"
	FormatPlaintext Format = "plaintext"
)

// ParseFormat resolves a --output flag value, defaulting to table.
// `--json` acts as a shorthand for FormatJSON.
func ParseFormat(output string, jsonShortcut bool) (Format, error) {
	if jsonShortcut {
		return FormatJSON, nil
	}
	if output == "" {
		return FormatTable, nil
	}
	switch f := Format(strings.ToLower(output)); f {
	case FormatTable, FormatJSON, FormatNDJSON, FormatYAML, FormatTSV, FormatPlaintext:
		return f, nil
	default:
		return "", fmt.Errorf("invalid --output %q: must be one of table, json, ndjson, yaml, tsv, plaintext", output)
	}
}

// ColumnDef describes one table column.
type ColumnDef struct {
	Header string
	// Path is a dotted JSON path (e.g. "user.email"), resolved against the
	// JSON-marshalled input value.
	Path string
}

// Render emits v in the requested format. For table/tsv, cols must be
// provided. fields (optional, comma-separated) performs client-side
// projection before rendering. jq (optional) is a final filter for
// JSON/NDJSON/YAML.
type RenderOpts struct {
	Format Format
	Cols   []ColumnDef
	Fields string
	JQ     string
	Pretty bool
}

// Render renders v according to opts.
func Render(w io.Writer, v any, opts RenderOpts) error {
	// JSON-normalise once; every format operates on this intermediate form.
	raw, err := toAny(v)
	if err != nil {
		return err
	}
	if opts.Fields != "" {
		raw = Project(raw, parseFields(opts.Fields))
	}
	if opts.JQ != "" {
		raw, err = applyJQ(raw, opts.JQ)
		if err != nil {
			return err
		}
	}
	switch opts.Format {
	case FormatJSON:
		return writeJSON(w, raw, opts.Pretty)
	case FormatNDJSON:
		return writeNDJSON(w, raw)
	case FormatYAML:
		return writeYAML(w, raw)
	case FormatTable:
		return writeTable(w, raw, opts.Cols)
	case FormatTSV:
		return writeTSV(w, raw, opts.Cols)
	case FormatPlaintext:
		return writePlaintext(w, raw)
	default:
		return fmt.Errorf("unknown format %q", opts.Format)
	}
}

func toAny(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	var out any
	if err := json.Unmarshal(b, &out); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return out, nil
}

func parseFields(s string) []string {
	out := []string{}
	for _, f := range strings.Split(s, ",") {
		if f = strings.TrimSpace(f); f != "" {
			out = append(out, f)
		}
	}
	return out
}

func writeJSON(w io.Writer, v any, pretty bool) error {
	enc := json.NewEncoder(w)
	if pretty {
		enc.SetIndent("", "  ")
	}
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func writeNDJSON(w io.Writer, v any) error {
	items, ok := v.([]any)
	if !ok {
		return writeJSON(w, v, false)
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	for _, item := range items {
		if err := enc.Encode(item); err != nil {
			return err
		}
	}
	return nil
}

func writeYAML(w io.Writer, v any) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		return err
	}
	return enc.Close()
}

func writePlaintext(w io.Writer, v any) error {
	switch t := v.(type) {
	case string:
		_, err := fmt.Fprintln(w, t)
		return err
	case []any:
		for _, item := range t {
			if s, ok := item.(string); ok {
				fmt.Fprintln(w, s)
			} else {
				b, _ := json.Marshal(item)
				fmt.Fprintln(w, string(b))
			}
		}
		return nil
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(b))
		return err
	}
}

func writeTable(w io.Writer, v any, cols []ColumnDef) error {
	rows, headers := rowsFromAny(v, cols)
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, r := range rows {
		for i, c := range r {
			if len(c) > widths[i] {
				widths[i] = len(c)
			}
		}
	}
	writeRow := func(cells []string) error {
		var b strings.Builder
		for i, c := range cells {
			if i > 0 {
				b.WriteString("  ")
			}
			if i == len(cells)-1 {
				b.WriteString(c)
			} else {
				b.WriteString(c)
				for n := len(c); n < widths[i]; n++ {
					b.WriteByte(' ')
				}
			}
		}
		b.WriteByte('\n')
		_, err := io.WriteString(w, b.String())
		return err
	}
	if err := writeRow(headers); err != nil {
		return err
	}
	for _, r := range rows {
		if err := writeRow(r); err != nil {
			return err
		}
	}
	return nil
}

func writeTSV(w io.Writer, v any, cols []ColumnDef) error {
	rows, headers := rowsFromAny(v, cols)
	cw := csv.NewWriter(w)
	cw.Comma = '\t'
	if err := cw.Write(headers); err != nil {
		return err
	}
	for _, r := range rows {
		if err := cw.Write(r); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func rowsFromAny(v any, cols []ColumnDef) (rows [][]string, headers []string) {
	headers = make([]string, len(cols))
	for i, c := range cols {
		headers[i] = c.Header
	}
	var items []any
	switch t := v.(type) {
	case []any:
		items = t
	default:
		items = []any{v}
	}
	for _, item := range items {
		row := make([]string, len(cols))
		for i, c := range cols {
			row[i] = stringify(pathLookup(item, c.Path))
		}
		rows = append(rows, row)
	}
	return
}

func pathLookup(v any, path string) any {
	if path == "" {
		return v
	}
	parts := strings.Split(path, ".")
	cur := v
	for _, p := range parts {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil
		}
		cur = m[p]
	}
	return cur
}

func stringify(v any) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%v", t)
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

func applyJQ(v any, expr string) (any, error) {
	q, err := gojq.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid --jq expression: %w", err)
	}
	iter := q.Run(v)
	var results []any
	for {
		r, ok := iter.Next()
		if !ok {
			break
		}
		if err, isErr := r.(error); isErr {
			return nil, fmt.Errorf("jq error: %w", err)
		}
		results = append(results, r)
	}
	if len(results) == 1 {
		return results[0], nil
	}
	return results, nil
}
