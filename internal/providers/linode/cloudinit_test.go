package linode

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"io"
	"testing"

	core "github.com/openclaw/crabbox/internal/cli"
)

func TestAcquireSendsGzipUserDataWithinLinodeLimit(t *testing.T) {
	api := &fakeLinodeAPI{}
	backend := newTestBackend(t, api)
	if _, err := backend.Acquire(context.Background(), core.AcquireRequest{Repo: core.Repo{Root: t.TempDir()}, RequestedSlug: "my-app"}); err != nil {
		t.Fatal(err)
	}
	raw, err := base64.StdEncoding.DecodeString(api.createRequests[0].Metadata.UserData)
	if err != nil {
		t.Fatal(err)
	}
	if len(raw) > 16*1024 {
		t.Fatalf("decoded user_data is %d bytes; Linode refuses over %d", len(raw), 16*1024)
	}
	gz, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("user_data is not gzip: %v", err)
	}
	cloudConfig, err := io.ReadAll(gz)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(cloudConfig, []byte("#cloud-config\n")) {
		t.Fatalf("gzip does not hold cloud-config: %.40q", cloudConfig)
	}
}
