package output

import "strings"

// Project keeps only the named fields (dotted paths supported) from v.
// Designed for top-level objects and arrays of objects. Unknown fields
// are silently dropped. Unsupported shapes (strings, numbers) are returned
// unchanged.
func Project(v any, fields []string) any {
	if len(fields) == 0 {
		return v
	}
	switch t := v.(type) {
	case []any:
		out := make([]any, len(t))
		for i, item := range t {
			out[i] = Project(item, fields)
		}
		return out
	case map[string]any:
		return projectMap(t, fields)
	default:
		return v
	}
}

func projectMap(m map[string]any, fields []string) map[string]any {
	out := map[string]any{}
	groups := map[string][]string{}
	for _, f := range fields {
		parts := strings.SplitN(f, ".", 2)
		head := parts[0]
		if len(parts) == 1 {
			groups[head] = append(groups[head], "")
		} else {
			groups[head] = append(groups[head], parts[1])
		}
	}
	for head, rest := range groups {
		val, ok := m[head]
		if !ok {
			continue
		}
		// If any rest is empty, user asked for the whole field; include it verbatim.
		wantAll := false
		sub := []string{}
		for _, r := range rest {
			if r == "" {
				wantAll = true
			} else {
				sub = append(sub, r)
			}
		}
		if wantAll || len(sub) == 0 {
			out[head] = val
			continue
		}
		out[head] = Project(val, sub)
	}
	return out
}
