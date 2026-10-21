// It is not possible to access global variables from Woodpecker,
// for example the server configuration. You must therefore parse
// the environment variables in your addon. The reason for this is
//  that the addon runs in a completely separate process.

package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"

	"go.woodpecker-ci.org/woodpecker/v3/server"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	forge_types "go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"golang.org/x/oauth2"
)

// Opts defines configuration options.
type Opts struct {
	URL               string // GitWeb URL
	OAuthClientID     string // Oauth2 client id.
	OAuthClientSecret string // Oauth2 client secret.
	OAuth2URL         string // User-facing Forgejo server url for OAuth2.
	OAuthHost         string // Public url for oauth if different from url.
	SkipVerify        bool   // Skip ssl verification.
}

// Gitolite Forge
type Gitolite struct {
	id                int64
	url               string
	oAuthClientID     string
	oAuthClientSecret string
	oAuth2URL         string
	oAuthHost         string
	skipVerify        bool
}

// New returns a Forge implementation that integrates with gitolite git
func New(id int64, opts Opts) (forge.Forge, error) {
	return &Gitolite{
		id:                id,
		url:               opts.URL,
		oAuthClientID:     opts.OAuthClientID,
		oAuthClientSecret: opts.OAuthClientSecret,
		oAuth2URL:         opts.OAuth2URL,
		oAuthHost:         opts.OAuthHost,
		skipVerify:        opts.SkipVerify,
	}, nil
}

// Name returns the string name of this driver.
func (g *Gitolite) Name() string {
	return "Gitolite"
}

// URL returns the root url of a configured forge.
func (g *Gitolite) URL() string {
	return g.url
}

func (g *Gitolite) oauth2Config(ctx context.Context) (*oauth2.Config, context.Context) {
	return &oauth2.Config{
			ClientID:     g.oAuthClientID,
			ClientSecret: g.oAuthClientSecret,
			Endpoint: oauth2.Endpoint{
				// TODO Must allow arbirtrary URLs
				AuthURL:  fmt.Sprintf("%s/api/oidc/authorization", g.url), // Authentication URL
				TokenURL: fmt.Sprintf("%s/api/oidc/token", g.url),         // Token URL
			},
			RedirectURL: fmt.Sprintf("%s/authorize", server.Config.Server.OAuthHost),
		},

		context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: g.skipVerify},
			Proxy:           http.ProxyFromEnvironment,
		}})
}

// Login authenticates a user via OAuth2.
//
// OAuth Flow:
//  1. Initial call with empty OAuthRequest.Code returns (nil, redirectURL, nil)
//  2. User authorizes at redirectURL
//  3. Second call with OAuthRequest.Code returns (User, redirectURL, nil)
//
// Returned User must contain: Login, Email, Avatar, AccessToken, RefreshToken, Expiry, ForgeRemoteID
func (g *Gitolite) Login(ctx context.Context, req *forge_types.OAuthRequest) (*model.User, string, error) {
	config, oauth2Ctx := g.oauth2Config(ctx)
	redirectURL := config.AuthCodeURL(req.State)

	// check the OAuth code
	if len(req.Code) == 0 {
		return nil, redirectURL, nil
	}

	token, err := config.Exchange(oauth2Ctx, req.Code)
	if err != nil {
		return nil, redirectURL, fmt.Errorf("error exchanging token: %w", err)
	}

	// TODO: this MUST be implemented (helper, forgejo uses forgejo.newClientToken)
	// client, err := c.newClientToken(ctx, token.AccessToken)
	// client, err := newClient(g.url, token.AccessToken, g.skipVerify)

	user := &model.User{
		Login:         "Forge User",
		Email:         "forge.user",
		ForgeRemoteID: model.ForgeRemoteID(fmt.Sprint("Jack")), // TODO: Needs to set to a real user
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		Expiry:        token.Expiry.UTC().Unix(),
	}

	return user, redirectURL, nil
}

func (g *Gitolite) Hook(ctx context.Context, r *http.Request) (*model.Repo, *model.Pipeline, error) {
	// TODO
	return nil, nil, types.ErrNotImplemented
}

func (g *Gitolite) Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error) {
	// TODO
	return nil, types.ErrNotImplemented
}

func (g *Gitolite) Repos(ctx context.Context, u *model.User, p *model.ListOptions) ([]*model.Repo, error) {
	// TODO
	return nil, types.ErrNotImplemented
}

func (g *Gitolite) File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, fileName string) ([]byte, error) {
	// TODO
	return nil, types.ErrNotImplemented
}

// NON Supported methods
func (g *Gitolite) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) Teams(ctx context.Context, u *model.User, p *model.ListOptions) ([]*model.Team, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) Dir(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, dirName string) ([]*types.FileMeta, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) Status(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, p *model.Workflow) error {
	return types.ErrNotImplemented
}
func (g *Gitolite) Activate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	return types.ErrNotImplemented
}
func (g *Gitolite) Deactivate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	return types.ErrNotImplemented
}
func (g *Gitolite) Branches(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]string, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) BranchHead(ctx context.Context, u *model.User, r *model.Repo, branch string) (*model.Commit, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) PullRequests(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]*model.PullRequest, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) OrgMembership(ctx context.Context, u *model.User, org string) (*model.OrgPerm, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) Org(ctx context.Context, u *model.User, org string) (*model.Org, error) {
	return nil, types.ErrNotImplemented
}
