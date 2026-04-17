// Package pagination provides offset pagination helpers.
package pagination

// Cursor tracks the current pagination position.
type Cursor struct {
	Skip     int
	Limit    int
	NextSkip *int
}

// NextPage returns a pointer to the next skip value when another page is
// available (i.e. returnedCount == limit), or nil when the last page has been
// reached. A limit of 0 means "no limit / unknown", and NextPage always
// returns nil in that case.
func NextPage(skip, limit, returnedCount int) *int {
	if limit <= 0 {
		return nil
	}
	if returnedCount < limit {
		return nil
	}
	next := skip + limit
	return &next
}

// NewCursor constructs a Cursor and computes NextSkip from the returned count.
func NewCursor(skip, limit, returnedCount int) Cursor {
	return Cursor{
		Skip:     skip,
		Limit:    limit,
		NextSkip: NextPage(skip, limit, returnedCount),
	}
}
