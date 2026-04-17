package output

// Envelope wraps a list response in the standard pagination envelope:
//
//	{"data": [...], "meta": {"limit": N, "skip": N, "next_skip": N|null}}
//
// Single-object responses must NOT be wrapped — callers are responsible for
// only calling Envelope on list results.
func Envelope(data any, limit, skip int, nextSkip *int) map[string]any {
	meta := map[string]any{
		"limit":     limit,
		"skip":      skip,
		"next_skip": nextSkip,
	}
	return map[string]any{
		"data": data,
		"meta": meta,
	}
}
