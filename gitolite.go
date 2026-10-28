package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.woodpecker-ci.org/woodpecker/v3/server"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	forge_types "go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"golang.org/x/oauth2"
)

// Opts defines configuration options.
type GitOpts struct {
	URL               string // GitWeb URL
	CloneURL          string // SSH clone URL
	OAuthClientID     string // Oauth2 client id.
	OAuthClientSecret string // Oauth2 client secret.
	OAuth2URL         string // OAuth2 provider URL to login
	SkipVerify        bool   // Skip ssl verification.
}

// Gitolite Forge interface
type Gitolite struct {
	url               string
	cloneurl          string
	oAuthClientID     string
	oAuthClientSecret string
	oAuth2URL         string
	skipVerify        bool
}

// New returns a Forge implementation that integrates with gitolite git
func New(opts GitOpts) forge.Forge {
	fmt.Fprintln(os.Stderr, "New Git Forge %s", opts.CloneURL)
	return &Gitolite{
		url:               opts.URL,
		cloneurl:          opts.CloneURL,
		oAuthClientID:     opts.OAuthClientID,
		oAuthClientSecret: opts.OAuthClientSecret,
		oAuth2URL:         opts.OAuth2URL,
		skipVerify:        opts.SkipVerify,
	}
}

// Name returns the string name of this driver.
func (g *Gitolite) Name() string {
	fmt.Fprintln(os.Stderr, "Name Gitolite Forge")
	return "Gitolite"
}

// URL returns the root url of a configured forge.
func (g *Gitolite) URL() string {
	fmt.Fprintln(os.Stderr, "URL Gitolite Forge %s", g.url)
	return g.url
}

func (g *Gitolite) oauth2Config(ctx context.Context) (*oauth2.Config, context.Context) {
	fmt.Fprintln(os.Stderr, "Gitolite OAuth2 Config: %s, %s", g.oAuth2URL, g.oAuthClientID)
	// OIDC Autodiscovery
	provider, err := oidc.NewProvider(ctx, g.oAuth2URL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failure initializing OIDC: %s", err.Error())
		return nil, nil
	}
	fmt.Fprintln(os.Stderr, "Got OIDC context")

	return &oauth2.Config{
			ClientID:     g.oAuthClientID,
			ClientSecret: g.oAuthClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  fmt.Sprintf("%s/authorize", server.Config.Server.OAuthHost),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
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
	fmt.Fprintln(os.Stderr, "Gitolite Login: %s", req.Code)
	config, oauth2Ctx := g.oauth2Config(ctx)
	redirectURL := config.AuthCodeURL(req.State)

	// check the OAuth code
	if len(req.Code) == 0 {
		fmt.Fprintln(os.Stderr, "OAuth code error %s", req.State)
		return nil, redirectURL, nil
	}

	token, err := config.Exchange(oauth2Ctx, req.Code)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Gitolite Auth Token Error: %s", err.Error())
		return nil, redirectURL, fmt.Errorf("error exchanging token: %w", err)
	}

	// TODO: this MUST be implemented (helper, forgejo uses forgejo.newClientToken)
	// client, err := c.newClientToken(ctx, token.AccessToken)
	// client, err := newClient(g.url, token.AccessToken, g.skipVerify)
	fmt.Fprintln(os.Stderr, "Got User Auth")
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
	fmt.Fprintln(os.Stderr, "Gitolite Hook")
	// TODO
	return nil, nil, types.ErrNotImplemented
}

func (g *Gitolite) Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error) {
	fmt.Fprintln(os.Stderr, "Gitolite Repo")
	// TODO
	return nil, types.ErrNotImplemented
}

func (g *Gitolite) Repos(ctx context.Context, u *model.User, p *model.ListOptions) ([]*model.Repo, error) {
	fmt.Fprintln(os.Stderr, "Gitolite Repos")
	// TODO
	return nil, types.ErrNotImplemented
}

func (g *Gitolite) File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, fileName string) ([]byte, error) {
	fmt.Fprintln(os.Stderr, "Gitolite File")
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
