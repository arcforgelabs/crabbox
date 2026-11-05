package linode

import (
	"encoding/base64"
	"fmt"

	core "github.com/openclaw/crabbox/internal/cli"
)

// linodeUserDataLimit is Linode's limit on decoded metadata user_data.
// https://techdocs.akamai.com/linode-api/reference/post-linode-instance
const linodeUserDataLimit = 16 * 1024

// linodeUserData gzips cloud-init, as Tencent Cloud does: the plain bootstrap is
// over Linode's limit, and cloud-init reads gzip user data as it is.
func linodeUserData(cfg core.Config, publicKey string) (string, error) {
	compressed, err := core.GzipCloudInit(core.CloudInitUserData(cfg, publicKey))
	if err != nil {
		return "", fmt.Errorf("gzip Linode user data: %w", err)
	}
	if len(compressed) > linodeUserDataLimit {
		return "", fmt.Errorf("Linode user-data limit is 16 KiB; compressed cloud-init is %d bytes", len(compressed))
	}
	return base64.StdEncoding.EncodeToString(compressed), nil
}
