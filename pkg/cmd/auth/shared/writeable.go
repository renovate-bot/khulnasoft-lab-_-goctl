package shared

import (
	"strings"

	"github.com/khulnasoft-lab/goctl/v2/internal/config"
)

func AuthTokenWriteable(authCfg *config.AuthConfig, hostname string) (string, bool) {
	token, src := authCfg.Token(hostname)
	return src, (token == "" || !strings.HasSuffix(src, "_TOKEN"))
}
