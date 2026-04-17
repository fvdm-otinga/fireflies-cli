# GraphQL Operation → CLI Command Map

All 35 Fireflies GraphQL operations and their corresponding CLI commands.

## Queries (17)

| GraphQL Operation | CLI Command | Notes |
|---|---|---|
| `user` (self) | `fireflies users whoami` | Returns the authenticated user |
| `users` | `fireflies users list` | List all workspace users |
| `user_groups` | `fireflies users groups` | List user groups |
| `channels` | `fireflies channels list` | List all channels |
| `channel` | `fireflies channels get <id>` | Get a specific channel by ID |
| `active_meetings` | `fireflies meetings active` | List currently live/active meetings |
| `transcripts` | `fireflies meetings list` | List meetings; supports `--since --until --limit --skip --transcript` |
| `transcript` | `fireflies meetings get <id>` | Get a single meeting; supports `--transcript --fields` |
| `transcript` (sentences) | `fireflies transcript text <id>` | Full plaintext with speaker attribution; supports `--since --until` |
| `bites` | `fireflies soundbites list` | List soundbites |
| `bite` | `fireflies soundbites get <id>` | Get a single soundbite by ID |
| `apps` | `fireflies apps list` | List Fireflies app integrations |
| `askfred_threads` | `fireflies askfred threads` | List AskFred AI threads |
| `askfred_thread` | `fireflies askfred thread <id>` | Get a specific AskFred thread |
| `contacts` | `fireflies contacts list` | List contacts |
| `analytics` | `fireflies analytics` | Fetch workspace analytics; supports `--since --until` |
| `live_action_items` | `fireflies live items <meeting-id>` | Fetch live action items for a meeting |
| `rule_executions_by_meeting` | `fireflies rules executions <meeting-id>` | Automation rule execution records |

## Mutations (18)

| GraphQL Operation | CLI Command | Notes |
|---|---|---|
| `uploadAudio` | `fireflies meetings upload <url>` | Upload from a public URL (single-step) |
| `createUploadUrl` | `fireflies meetings upload <file>` | Step 1 of local file upload |
| `confirmUpload` | `fireflies meetings upload <file>` | Step 2 of local file upload (auto-chained) |
| `updateMeetingTitle` | `fireflies meetings update title <id> <title>` | Rename a meeting |
| `updateMeetingPrivacy` | `fireflies meetings update privacy <id> <level>` | Change meeting privacy level |
| `updateMeetingState` | `fireflies meetings update state <id> <state>` | Change meeting processing state |
| `updateMeetingChannel` | `fireflies meetings move <id> --channel <cid>` | Move meeting to a channel |
| `shareMeeting` | `fireflies meetings share <id> --email <e>...` | Share with one or more email addresses; rate: 10/hr |
| `revokeSharedMeetingAccess` | `fireflies meetings revoke <id> --email <e>` | Revoke access for an email [destructive] |
| `deleteTranscript` | `fireflies meetings delete <id>` | Delete a meeting [destructive]; rate: 10/min |
| `createBite` | `fireflies soundbites create --meeting <id> --start <t> --end <t>` | Create a soundbite clip |
| `addToLiveMeeting` | `fireflies live add` | Add participant to live meeting; rate: 3/20min |
| `createLiveSoundbite` | `fireflies live soundbite` | Create a soundbite during a live meeting |
| `createLiveActionItem` | `fireflies live action-item` | Create an action item during a live meeting |
| `createAskFredThread` | `fireflies askfred ask --meeting <id> <question>` | Start a new AskFred AI thread |
| `continueAskFredThread` | `fireflies askfred continue <thread-id> <question>` | Continue an existing AskFred thread |
| `deleteAskFredThread` | `fireflies askfred delete <thread-id>` | Delete an AskFred thread [destructive] |
| `setUserRole` | `fireflies users set-role <uid> <role>` | Change a user's workspace role [destructive] |

## Coverage

- Queries: 17/17 (plus `transcript` is exposed via two commands for different use-cases)
- Mutations: 18/18 (with `createUploadUrl`+`confirmUpload` auto-chained in `meetings upload`)
- **Total: 35/35 GraphQL operations covered**

## Other surfaces

| Surface | CLI Command | Notes |
|---|---|---|
| Socket.IO realtime | `fireflies realtime tail <meeting-id>` | Streams transcript events as JSON or plaintext |
| Webhooks V1+V2 | `fireflies webhooks serve` | HTTP receiver with HMAC-SHA256 signature verification |
