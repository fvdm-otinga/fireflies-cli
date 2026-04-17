## fireflies meetings

Query and manage Fireflies meetings (transcripts)

### Options

```
  -h, --help   help for meetings
```

### Options inherited from parent commands

```
      --dry-run             Print the GraphQL operation without executing
      --fields string       Comma-separated top-level fields to keep (client-side projection)
      --jq string           Post-process output via a gojq expression
      --json                Shortcut for --output json
      --limit int           Page size (0 = API default, max 50 for transcripts)
      --output string       Output format: table|json|ndjson|yaml|tsv|plaintext
      --profile string      Config profile to use
      --since string        Lower bound (RFC3339 or relative like 7d)
      --skip int            Offset pagination cursor
      --transcript string   Transcript depth: none|preview|full
      --until string        Upper bound (RFC3339)
      --yes                 Bypass confirmation prompts for destructive operations
```

### SEE ALSO

* [fireflies](fireflies.md)	 - Fireflies.ai CLI (token-efficient wrapper for the GraphQL API)
* [fireflies meetings active](fireflies_meetings_active.md)	 - List active/live meetings (GraphQL: active_meetings)
* [fireflies meetings delete](fireflies_meetings_delete.md)	 - Delete a meeting/transcript (GraphQL: deleteTranscript, rate: 10/min) [destructive]
* [fireflies meetings get](fireflies_meetings_get.md)	 - Get a single meeting/transcript by ID (GraphQL: transcript)
* [fireflies meetings list](fireflies_meetings_list.md)	 - List meetings/transcripts (GraphQL: transcripts)
* [fireflies meetings move](fireflies_meetings_move.md)	 - Move a meeting to a channel (GraphQL: updateMeetingChannel)
* [fireflies meetings revoke](fireflies_meetings_revoke.md)	 - Revoke shared meeting access for an email (GraphQL: revokeSharedMeetingAccess) [destructive]
* [fireflies meetings share](fireflies_meetings_share.md)	 - Share a meeting with one or more email addresses (GraphQL: shareMeeting, rate: 10/hr)
* [fireflies meetings update](fireflies_meetings_update.md)	 - Update meeting metadata
* [fireflies meetings upload](fireflies_meetings_upload.md)	 - Upload an audio/video file or URL to Fireflies (GraphQL: uploadAudio or createUploadUrl+confirmUpload)

