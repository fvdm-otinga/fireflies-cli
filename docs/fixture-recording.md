# Fixture Recording Guide

go-vcr cassettes in `testdata/fixtures/` allow CI to run tests without hitting the live Fireflies API.

---

## How Cassettes Work

1. **Record mode** (`TEST_RECORD=1`): the test makes a real HTTP request and saves the interaction to a YAML cassette.
2. **Replay mode** (default): the test replays the cassette without network access. CI always runs in replay mode.
3. **Scrub**: before committing, run `scripts/scrub-fixtures.sh` to redact Bearer tokens, email addresses, and UUIDs.

---

## Recording a New Cassette

```sh
# Set a live API key
export FIREFLIES_API_KEY=ff_xxxxxxxxxxxxxxxx

# Record (replaces existing cassette if present)
TEST_RECORD=1 go test ./internal/client/... -run TestWhoami_Replay -v

# Scrub secrets from the newly recorded cassette
make scrub

# Verify scrub was clean
make scrub-check

# Commit
git add testdata/fixtures/whoami_happy.yaml
git commit -m "chore: refresh whoami_happy.yaml cassette"
```

---

## Cassette File Format

Cassettes use go-vcr v4 format (version: 2). Each file lives at:

```
testdata/fixtures/<command>_happy.yaml
```

Example cassette structure:

```yaml
---
version: 2
interactions:
    - id: 0
      request:
        proto: HTTP/1.1
        host: api.fireflies.ai
        headers:
            Authorization:
                - "[SCRUBBED]"
            Content-Type:
                - application/json
        body: '{"query":"...", "operationName":"..."}'
        url: https://api.fireflies.ai/graphql
        method: POST
      response:
        proto: HTTP/2.0
        content_length: -1
        uncompressed: true
        headers:
            Content-Type:
                - application/json; charset=utf-8
        body: '{"data":{...scrubbed...}}'
        status: 200 OK
        code: 200
```

---

## Commands Still Needing Cassettes (post-v1.0.0)

For v1.0.0, only `whoami_happy.yaml` is recorded. All others are planned post-v1. Priority order:

| Cassette | GraphQL op | Command |
|---|---|---|
| `meetings_list_happy.yaml` | `transcripts` | `meetings list` |
| `meetings_get_happy.yaml` | `transcript` | `meetings get <id>` |
| `users_list_happy.yaml` | `users` | `users list` |
| `channels_list_happy.yaml` | `channels` | `channels list` |
| `analytics_happy.yaml` | `analytics` | `analytics` |
| `soundbites_list_happy.yaml` | `bites` | `soundbites list` |
| `contacts_list_happy.yaml` | `contacts` | `contacts list` |
| `askfred_threads_happy.yaml` | `askfred_threads` | `askfred threads` |

For mutations, record with a dedicated test account and use `--dry-run` cassettes (zero HTTP) where possible.

---

## Scrub Patterns

`scripts/scrub-fixtures.sh` redacts these patterns (ERE):

| Pattern | What it matches |
|---|---|
| `Bearer [A-Za-z0-9._/+=-]+` | API keys in Authorization headers |
| `[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}` | Email addresses |
| `[0-9a-fA-F]{8}-...-[0-9a-fA-F]{12}` | UUIDs |
| `FIREFLIES_API_KEY=[^ ]+` | Env var values |

After scrubbing, replace Authorization values with `[SCRUBBED]` and email fields with `[scrubbed]` (no `@`) so the scrub check does not false-positive on the REDACTED placeholders themselves.

---

## Matcher Configuration

The replay test uses a URL+method-only matcher (ignoring request body and auth headers):

```go
urlMatcher := cassette.MatcherFunc(func(r *http.Request, i cassette.Request) bool {
    return r.Method == i.Method && r.URL.String() == i.URL
})
```

This makes cassettes robust against minor request-body whitespace differences and allows replay with any API key value.
