package meetings

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/fvdm-otinga/fireflies-cli/internal/client"
	"github.com/fvdm-otinga/fireflies-cli/internal/config"
	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
	"github.com/fvdm-otinga/fireflies-cli/internal/flags"
	ffgql "github.com/fvdm-otinga/fireflies-cli/internal/graphql"
	"github.com/fvdm-otinga/fireflies-cli/internal/netguard"
	"github.com/fvdm-otinga/fireflies-cli/internal/output"
)

// newUploadCmd returns `fireflies meetings upload <file-or-url>`.
// GraphQL: UploadAudio (URL) or CreateUploadUrl+ConfirmUpload (local file).
func newUploadCmd() *cobra.Command {
	var (
		title           string
		attendees       []string
		customLanguage  string
		saveVideo       bool
		webhook         string
		clientRefID     string
	)

	cmd := &cobra.Command{
		Use:   "upload <file-or-url>",
		Short: "Upload an audio/video file or URL to Fireflies (GraphQL: uploadAudio or createUploadUrl+confirmUpload)",
		Long: `Upload audio or video to Fireflies for transcription.

  URL input:  uses the uploadAudio mutation (single-step).
  File input: uses createUploadUrl → S3 PUT → confirmUpload (two-step).

Required: --title`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := flags.FromCmd(cmd)
			target := args[0]

			if title == "" {
				return ferr.Usage("--title is required")
			}

			// Detect if target is a local file or a URL.
			isLocal := false
			if _, err := os.Stat(target); err == nil {
				isLocal = true
			}

			if sh.DryRun {
				if isLocal {
					// Two-step: createUploadUrl + confirmUpload
					attendeeInputs := buildAttendeeInputs(attendees)
					contentType := "audio/mpeg"
					fileSize := 0
					vars1, _ := json.MarshalIndent(map[string]any{
						"input": map[string]any{
							"title":           title,
							"content_type":    contentType,
							"file_size":       fileSize,
							"custom_language": customLanguage,
							"attendees":       attendeeInputs,
						},
					}, "", "  ")
					_, _ = fmt.Fprintf(os.Stdout, "mutation CreateUploadUrl($input: CreateUploadUrlInput!) {\n  createUploadUrl(input: $input) { upload_url meeting_id expires_at }\n}\n")
					_, _ = fmt.Fprintf(os.Stdout, "%s\n", vars1)
					_, _ = fmt.Fprintf(os.Stdout, "# then: S3 PUT <upload_url>\n")
					vars2, _ := json.MarshalIndent(map[string]any{
						"input": map[string]any{"meeting_id": "<meeting_id from step 1>"},
					}, "", "  ")
					_, _ = fmt.Fprintf(os.Stdout, "mutation ConfirmUpload($input: ConfirmUploadInput!) {\n  confirmUpload(input: $input) { success meeting_id message }\n}\n")
					_, _ = fmt.Fprintf(os.Stdout, "%s\n", vars2)
				} else {
					// Single-step: uploadAudio
					attendeeInputs := buildAttendeeInputs(attendees)
					vars, _ := json.MarshalIndent(map[string]any{
						"input": map[string]any{
							"url":                 target,
							"title":               title,
							"attendees":           attendeeInputs,
							"custom_language":     customLanguage,
							"save_video":          saveVideo,
							"webhook":             webhook,
							"client_reference_id": clientRefID,
						},
					}, "", "  ")
					_, _ = fmt.Fprintf(os.Stdout, "mutation UploadAudio($input: AudioUploadInput) {\n  uploadAudio(input: $input) { success title message }\n}\n")
					_, _ = fmt.Fprintf(os.Stdout, "%s\n", vars)
				}
				return nil
			}

			prof, err := config.New().Profile(sh.Profile)
			if err != nil {
				return err
			}
			c := client.New(prof)

			if isLocal {
				return uploadLocalFile(cmd.Context(), c, target, title, attendees, customLanguage, sh, cmd.OutOrStdout())
			}
			return uploadURL(cmd.Context(), c, target, title, attendees, customLanguage, saveVideo, webhook, clientRefID, sh, cmd.OutOrStdout())
		},
	}

	flags.Bind(cmd)
	cmd.Flags().StringVar(&title, "title", "", "Title for the uploaded meeting (required)")
	cmd.Flags().StringSliceVar(&attendees, "attendees", nil, "Attendee emails (comma-separated)")
	cmd.Flags().StringVar(&customLanguage, "custom-language", "", "Custom language code for transcription")
	cmd.Flags().BoolVar(&saveVideo, "save-video", false, "Save video in addition to audio")
	cmd.Flags().StringVar(&webhook, "webhook", "", "Webhook URL to notify on completion")
	cmd.Flags().StringVar(&clientRefID, "client-reference-id", "", "Client reference ID")

	return cmd
}

func buildAttendeeInputs(emails []string) []map[string]any {
	if len(emails) == 0 {
		return nil
	}
	out := make([]map[string]any, 0, len(emails))
	for _, e := range emails {
		out = append(out, map[string]any{"email": strings.TrimSpace(e)})
	}
	return out
}

func toAttendeeInput(emails []string) []*ffgql.AttendeeInput {
	if len(emails) == 0 {
		return nil
	}
	out := make([]*ffgql.AttendeeInput, 0, len(emails))
	for _, e := range emails {
		email := strings.TrimSpace(e)
		out = append(out, &ffgql.AttendeeInput{Email: &email})
	}
	return out
}

func uploadURL(ctx context.Context, c *client.Client, url, title string, attendees []string, customLang string, saveVideo bool, webhook, clientRefID string, sh *flags.Shared, w io.Writer) error {
	input := &ffgql.AudioUploadInput{
		Url:   url,
		Title: &title,
	}
	if customLang != "" {
		input.Custom_language = &customLang
	}
	if saveVideo {
		input.Save_video = &saveVideo
	}
	if webhook != "" {
		input.Webhook = &webhook
	}
	if clientRefID != "" {
		input.Client_reference_id = &clientRefID
	}
	input.Attendees = toAttendeeInput(attendees)

	resp, err := ffgql.UploadAudio(ctx, c, input)
	if err != nil {
		return ferr.FromGraphQLError(err)
	}

	f, err := output.ParseFormat(sh.Output, sh.JSON)
	if err != nil {
		return ferr.Usage(err.Error())
	}
	return output.Render(w, resp.UploadAudio, output.RenderOpts{
		Format: f,
		Cols: []output.ColumnDef{
			{Header: "SUCCESS", Path: "success"},
			{Header: "TITLE", Path: "title"},
			{Header: "MESSAGE", Path: "message"},
		},
		Fields: sh.Fields,
		JQ:     sh.JQ,
		Pretty: sh.JSON,
	})
}

func uploadLocalFile(ctx context.Context, c *client.Client, filePath, title string, attendees []string, customLang string, sh *flags.Shared, w io.Writer) error {
	fi, err := os.Stat(filePath)
	if err != nil {
		return ferr.General(fmt.Sprintf("cannot stat file: %v", err))
	}

	contentType := "audio/mpeg"
	if strings.HasSuffix(strings.ToLower(filePath), ".mp4") || strings.HasSuffix(strings.ToLower(filePath), ".mov") {
		contentType = "video/mp4"
	} else if strings.HasSuffix(strings.ToLower(filePath), ".m4a") {
		contentType = "audio/mp4"
	} else if strings.HasSuffix(strings.ToLower(filePath), ".wav") {
		contentType = "audio/wav"
	} else if strings.HasSuffix(strings.ToLower(filePath), ".ogg") {
		contentType = "audio/ogg"
	}

	fileSize := int(fi.Size())
	input := &ffgql.CreateUploadUrlInput{
		Title:        &title,
		Content_type: contentType,
		File_size:    fileSize,
	}
	if customLang != "" {
		input.Custom_language = &customLang
	}
	input.Attendees = toAttendeeInput(attendees)

	// Step 1: Get signed upload URL.
	urlResp, err := ffgql.CreateUploadUrl(ctx, c, input)
	if err != nil {
		return ferr.FromGraphQLError(err)
	}
	uploadURL := urlResp.CreateUploadUrl.Upload_url
	meetingID := urlResp.CreateUploadUrl.Meeting_id

	// SSRF guard: validate the signed upload URL before we trust it.
	if _, err := netguard.ValidateUploadURL(uploadURL); err != nil {
		return ferr.General(err.Error())
	}
	_, _ = fmt.Fprintf(w, "Uploading %s (%d bytes) to meeting %s...\n", filePath, fileSize, meetingID)

	// Step 2: PUT file to S3.
	f, err := os.Open(filePath)
	if err != nil {
		return ferr.General(fmt.Sprintf("open file: %v", err))
	}
	defer f.Close() //nolint:errcheck // read-only open; close error not actionable

	putReq, err := http.NewRequestWithContext(ctx, http.MethodPut, uploadURL, &progressReader{r: f, total: int64(fileSize), w: w})
	if err != nil {
		return ferr.General(fmt.Sprintf("build S3 PUT: %v", err))
	}
	putReq.ContentLength = int64(fileSize)
	putReq.Header.Set("Content-Type", contentType)

	// Hardened client for the S3 PUT: enforce TLS 1.2+, a 30-minute ceiling,
	// and refuse to follow redirects (a redirect could point at an internal
	// host that slipped past ValidateUploadURL).
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	s3Client := &http.Client{
		Timeout:   30 * time.Minute,
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	putResp, err := s3Client.Do(putReq)
	if err != nil {
		return ferr.General(fmt.Sprintf("S3 PUT failed: %v", err))
	}
	defer putResp.Body.Close() //nolint:errcheck // response body close
	if putResp.StatusCode < 200 || putResp.StatusCode >= 300 {
		body, _ := io.ReadAll(putResp.Body)
		return ferr.General(fmt.Sprintf("S3 PUT returned %d: %s", putResp.StatusCode, string(body)))
	}
	_, _ = fmt.Fprintln(w, " done.")

	// Step 3: Confirm upload.
	confirmResp, err := ffgql.ConfirmUpload(ctx, c, &ffgql.ConfirmUploadInput{Meeting_id: meetingID})
	if err != nil {
		return ferr.FromGraphQLError(err)
	}

	fmtOpt, err := output.ParseFormat(sh.Output, sh.JSON)
	if err != nil {
		return ferr.Usage(err.Error())
	}
	return output.Render(w, confirmResp.ConfirmUpload, output.RenderOpts{
		Format: fmtOpt,
		Cols: []output.ColumnDef{
			{Header: "SUCCESS", Path: "success"},
			{Header: "MEETING_ID", Path: "meeting_id"},
			{Header: "MESSAGE", Path: "message"},
		},
		Fields: sh.Fields,
		JQ:     sh.JQ,
		Pretty: sh.JSON,
	})
}

// progressReader wraps an io.Reader and prints progress.
type progressReader struct {
	r     io.Reader
	total int64
	read  int64
	w     io.Writer
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.read += int64(n)
	if p.total > 0 {
		pct := float64(p.read) / float64(p.total) * 100
		_, _ = fmt.Fprintf(p.w, "\r  %.0f%%", pct)
	}
	return n, err
}
