package main

import (
	"os"

	"go.woodpecker-ci.org/woodpecker/v3/server/forge/addon"
)

// Create new forge from environment settings
func main() {
	// Get config environment
	opts := GitOpts{
		URL:               os.Getenv("GITWEB_URL"),
		CloneURL:          os.Getenv("GITSSH_URL"),
		OAuthClientID:     os.Getenv("GITAUTH_ID"),
		OAuthClientSecret: os.Getenv("GITAUTH_SECRET"),
		OAuth2URL:         os.Getenv("GITAUTH_URL"),
		SkipVerify:        true,
	}

	// Serve Addon
	addon.Serve(New(opts))
}
