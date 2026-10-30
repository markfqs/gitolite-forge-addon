package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"go.woodpecker-ci.org/woodpecker/v3/server"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge"
	"go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	forge_types "go.woodpecker-ci.org/woodpecker/v3/server/forge/types"
	"go.woodpecker-ci.org/woodpecker/v3/server/model"
	"golang.org/x/oauth2"
)

// GitOpts defines configuration options.
type GitOpts struct {
	URL               string // GitWeb URL
	CloneURL          string // SSH clone URL
	OAuthClientID     string // Oauth2 client id.
	OAuthClientSecret string // Oauth2 client secret.
	OAuth2URL         string // OAuth2 provider URL to login
	SkipVerify        bool   // Skip ssl verification.
}

// Gitolite Forge
type Gitolite struct {
	url               string
	cloneurl          string
	oAuthClientID     string
	oAuthClientSecret string
	oAuth2URL         string
	skipVerify        bool
}

// Forge implementation that integrates with gitolite git
func NewGitolite(opts GitOpts) forge.Forge {
	fmt.Fprintf(os.Stderr, "New Git Forge %s\n", opts.CloneURL)
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
	fmt.Fprintf(os.Stderr, "Name Gitolite Forge")
	return "Gitolite"
}

// URL returns the root url of a configured forge.
func (g *Gitolite) URL() string {
	fmt.Fprintf(os.Stderr, "URL Gitolite Forge %s\n", g.url)
	return g.url
}

// Return oauth config using OpenID autodiscovery
func (g *Gitolite) oauth2Config(ctx context.Context) (*oauth2.Config, context.Context) {
	fmt.Fprintf(os.Stderr, "Gitolite OAuth2 Config: %s\n, %s\n", g.oAuth2URL, g.oAuthClientID)

	// OIDC Autodiscovery
	provider, err := oidc.NewProvider(ctx, g.oAuth2URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failure discovering OAuth2 provider: %s\n", err.Error())
		return nil, nil
	}
	fmt.Fprintf(os.Stderr, "Got OAuth2 provider")

	return &oauth2.Config{
			ClientID:     g.oAuthClientID,
			ClientSecret: g.oAuthClientSecret,
			Endpoint:     provider.Endpoint(),
			Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
			RedirectURL:  fmt.Sprintf("%s\n/authorize", server.Config.Server.OAuthHost), // CHECK: Not sure if can access server.Config
		},

		context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: g.skipVerify},
			Proxy:           http.ProxyFromEnvironment,
		}})
}

// Login authenticates a user via OAuth2.
// Flow into passes, empty cal, auth, return user identity
func (g *Gitolite) Login(ctx context.Context, req *forge_types.OAuthRequest) (*model.User, string, error) {
	fmt.Fprintf(os.Stderr, "Gitolite Login: %s\n", req.Code)

	config, oauth2Ctx := g.oauth2Config(ctx)
	redirectURL := config.AuthCodeURL(req.State)

	// check the OAuth code
	if len(req.Code) == 0 {
		fmt.Fprintf(os.Stderr, "OAuth First pass")
		return nil, redirectURL, nil
	}

	token, err := config.Exchange(oauth2Ctx, req.Code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Gitolite Auth Token Error: %s\n", err.Error())
		return nil, redirectURL, fmt.Errorf("error exchanging token: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Got User Auth")
	user := &model.User{
		Login:         "dummy",
		Email:         "dummy@forge.user",
		ForgeRemoteID: "dummy",
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		Expiry:        time.Now().Add(time.Hour).Unix(),
	}

	return user, redirectURL, nil
}

/*
REQUIRED METHODS : Mimic/Mock TODO:Implement
*/
func (g *Gitolite) Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error) {
	fmt.Fprintf(os.Stderr, "Gitolite Repo")
	reponame := name
	return &model.Repo{
		ForgeRemoteID: remoteID,
		Name:          reponame,
		FullName:      reponame + " repository",
		Owner:         owner,
		ForgeURL:      g.url + "/" + reponame,
		Clone:         g.url + "/" + reponame,
		CloneSSH:      g.cloneurl + "/" + reponame,
		Branch:        "master",
		PREnabled:     false,
		Perm: &model.Perm{
			Pull:  true,
			Push:  false,
			Admin: false,
		},
		IsSCMPrivate: true,
	}, nil
}

func (g *Gitolite) Hook(ctx context.Context, r *http.Request) (*model.Repo, *model.Pipeline, error) {
	fmt.Fprintf(os.Stderr, "Gitolite Hook")
	// TODO
	return nil, nil, &types.ErrIgnoreEvent{}
}

func (g *Gitolite) Repos(ctx context.Context, u *model.User, p *model.ListOptions) ([]*model.Repo, error) {
	fmt.Fprintf(os.Stderr, "Gitolite Repos")
	// TODO

	var result []*model.Repo
	repo, _ := g.Repo(ctx, u, "1", "dummy", "dummy")
	result = append(result, repo)
	return result, nil
}

func (g *Gitolite) File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, fileName string) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "Gitolite File")
	// TODO
	return nil, nil
}

func (g *Gitolite) Status(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, p *model.Workflow) error {
	return nil
}

func (g *Gitolite) Netrc(u *model.User, r *model.Repo) (*model.Netrc, error) {
	return &model.Netrc{
		Machine:  "localhost",
		Login:    "dummy",
		Password: "nopass",
	}, nil
}

func (g *Gitolite) Activate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	return nil
}

func (g *Gitolite) Deactivate(ctx context.Context, u *model.User, r *model.Repo, link string) error {
	return nil
}

func (g *Gitolite) BranchHead(ctx context.Context, u *model.User, r *model.Repo, branch string) (*model.Commit, error) {
	return &model.Commit{
		SHA:      "",
		ForgeURL: r.ForgeURL,
	}, nil
}

func (g *Gitolite) Org(ctx context.Context, u *model.User, org string) (*model.Org, error) {
	return &model.Org{
		Name:   g.Name(),
		IsUser: true,
	}, nil
}

// NON Supported optional methods
func (g *Gitolite) Teams(ctx context.Context, u *model.User, p *model.ListOptions) ([]*model.Team, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) Dir(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, dirName string) ([]*types.FileMeta, error) {
	return nil, types.ErrNotImplemented
}

func (g *Gitolite) Branches(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]string, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) PullRequests(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]*model.PullRequest, error) {
	return nil, types.ErrNotImplemented
}
func (g *Gitolite) OrgMembership(ctx context.Context, u *model.User, org string) (*model.OrgPerm, error) {
	return nil, types.ErrNotImplemented
}
