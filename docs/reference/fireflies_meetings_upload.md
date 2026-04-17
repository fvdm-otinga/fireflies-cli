## fireflies meetings upload

Upload an audio/video file or URL to Fireflies (GraphQL: uploadAudio or createUploadUrl+confirmUpload)

### Synopsis

Upload audio or video to Fireflies for transcription.

  URL input:  uses the uploadAudio mutation (single-step).
  File input: uses createUploadUrl → S3 PUT → confirmUpload (two-step).

Required: --title

```
fireflies meetings upload <file-or-url> [flags]
```

### Options

```
      --attendees strings            Attendee emails (comma-separated)
      --client-reference-id string   Client reference ID
      --custom-language string       Custom language code for transcription
      --dry-run                      Print the GraphQL operation without executing
      --fields string                Comma-separated top-level fields to keep (client-side projection)
  -h, --help                         help for upload
      --jq string                    Post-process output via a gojq expression
      --json                         Shortcut for --output json
      --limit int                    Page size (0 = API default, max 50 for transcripts)
      --output string                Output format: table|json|ndjson|yaml|tsv|plaintext
      --profile string               Config profile to use
      --save-video                   Save video in addition to audio
      --since string                 Lower bound (RFC3339 or relative like 7d)
      --skip int                     Offset pagination cursor
      --title string                 Title for the uploaded meeting (required)
      --transcript string            Transcript depth: none|preview|full
      --until string                 Upper bound (RFC3339)
      --webhook string               Webhook URL to notify on completion
      --yes                          Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies meetings](fireflies_meetings.md)	 - Query and manage Fireflies meetings (transcripts)

