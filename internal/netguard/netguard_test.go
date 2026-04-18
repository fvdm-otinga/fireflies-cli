package netguard

import (
	"strings"
	"testing"
)

func TestValidateUploadURL(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		wantErr bool
	}{
		{
			name: "https public S3 host ok",
			raw:  "https://fireflies-uploads.s3.amazonaws.com/abc?sig=xyz",
		},
		{
			name:    "http rejected",
			raw:     "http://fireflies-uploads.s3.amazonaws.com/abc",
			wantErr: true,
		},
		{
			name:    "localhost rejected",
			raw:     "https://localhost/abc",
			wantErr: true,
		},
		{
			name:    "AWS metadata IP rejected",
			raw:     "https://169.254.169.254/latest/meta-data/",
			wantErr: true,
		},
		{
			name:    "RFC1918 10/8 rejected",
			raw:     "https://10.0.0.1/abc",
			wantErr: true,
		},
		{
			name:    "RFC1918 192.168/16 rejected",
			raw:     "https://192.168.1.1/abc",
			wantErr: true,
		},
		{
			name:    "loopback 127/8 rejected",
			raw:     "https://127.0.0.1/abc",
			wantErr: true,
		},
		{
			name:    "IPv6 loopback rejected",
			raw:     "https://[::1]/abc",
			wantErr: true,
		},
		{
			name:    "link-local 169.254/16 rejected",
			raw:     "https://169.254.10.1/abc",
			wantErr: true,
		},
		{
			name:    "CGNAT 100.64/10 rejected",
			raw:     "https://100.64.1.2/abc",
			wantErr: true,
		},
		{
			name:    "metadata.google.internal rejected",
			raw:     "https://metadata.google.internal/computeMetadata/v1/",
			wantErr: true,
		},
		{
			name:    "metadata.goog rejected",
			raw:     "https://metadata.goog/",
			wantErr: true,
		},
		{
			name:    "ftp scheme rejected",
			raw:     "ftp://example.com/abc",
			wantErr: true,
		},
		{
			name:    "unparseable URL rejected",
			raw:     "://bad",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			u, err := ValidateUploadURL(tc.raw)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (url=%v)", u)
				}
				if !strings.Contains(err.Error(), "invalid upload URL scheme") {
					t.Fatalf("error message should mention 'invalid upload URL scheme', got: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if u == nil {
				t.Fatalf("expected non-nil URL on success")
			}
		})
	}
}
