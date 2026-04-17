package dynamic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Khan/genqlient/graphql"
)

// TranscriptResult holds a dynamically decoded Transcript object. All fields
// are json.RawMessage so callers can unmarshal only what they need.
type TranscriptResult map[string]json.RawMessage

// TranscriptsListResult is a list of dynamically decoded transcripts.
type TranscriptsListResult []map[string]json.RawMessage

// ExecuteSingleTranscript runs a DynamicTranscript query with the given fields
// and returns the raw decoded object.
func ExecuteSingleTranscript(ctx context.Context, c graphql.Client, id string, fields []TranscriptField) (map[string]any, error) {
	query := BuildSingleTranscriptQuery(fields)
	vars := map[string]any{"id": id}

	type resp struct {
		Transcript map[string]any `json:"transcript"`
	}
	var out resp
	r := &graphql.Response{Data: &out}
	err := c.MakeRequest(ctx, &graphql.Request{
		OpName:    "DynamicTranscript",
		Query:     query,
		Variables: vars,
	}, r)
	if err != nil {
		return nil, fmt.Errorf("dynamic transcript query: %w", err)
	}
	return out.Transcript, nil
}

// ExecuteTranscriptsList runs a DynamicTranscripts query and returns the raw slice.
func ExecuteTranscriptsList(ctx context.Context, c graphql.Client, vars map[string]any, fields []TranscriptField) ([]any, error) {
	query := BuildTranscriptsListQuery(fields)

	type resp struct {
		Transcripts []any `json:"transcripts"`
	}
	var out resp
	r := &graphql.Response{Data: &out}
	err := c.MakeRequest(ctx, &graphql.Request{
		OpName:    "DynamicTranscripts",
		Query:     query,
		Variables: vars,
	}, r)
	if err != nil {
		return nil, fmt.Errorf("dynamic transcripts query: %w", err)
	}
	return out.Transcripts, nil
}
