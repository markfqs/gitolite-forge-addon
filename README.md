Just compile/test:

go mod init gitolite-forge && go mod tidy && go build -o gitolite ./gitolite

Login
 it also supplies the user identity and tokens that Woodpecker needs for subsequent operations

# IDEAS

Name it gitolite
Forge-Addon will use a fixed provided SSH key to access the repository

environment:
  FORGE_SSH_HOST: forge.internal
  FORGE_SSH_USER: woodpecker
  FORGE_SSH_KEY_PATH: /run/secrets/forge_ssh_key
  FORGE_SSH_KNOWN_HOSTS: /run/secrets/forge_known_hosts

Probablay would need to package separate helper code for SSH interactions:
- Repository listing by user
- Repository metadata by user

Implement:
- Login -> OAuth
- Hook
- File -> SSH
- Repo -> SSH
- Repos -> SSH

How to mock them?
There is an OLD go package for gitolite:
https://github.com/VonC/gogitolite
https://pkg.go.dev/github.com/VonC/gogitolite/gitolite

For SSH need:
- Known hosts
- SSH user and key
- SSH server and port

	// Repo fetches a single repository.
	//
	// Lookup Strategy:
	// - Prefer lookup by remoteID (forge's internal ID) if provided (more reliable as repos can be renamed)
	// - Fallback to owner/name if remoteID empty
	//
	// Must verify user has at least read access.
	// Caller must make sure ForgeID is set.
	Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error)

	// Repos fetches all repositories accessible to the user.
	// Should include user's permission level in Repo.Perm.
	// Should support pagination via ListOptions.
	// Caller must make sure ForgeID is set.
	Repos(ctx context.Context, u *model.User, p *model.ListOptions) ([]*model.Repo, error)

	// File fetches a single file at a specific commit.
	// Primary method for retrieving pipeline configuration files.
	// Must fetch at specific commit (b.Commit), not branch head.
	File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, fileName string) ([]byte, error)

	// Hook parses incoming webhook and returns pipeline data.
	//
	// Webhook Processing Flow:
	//  1. HTTP request arrives at /api/hook with forge-specific format
	//  2. Webhook token verified against repo.Hash
	//  3. Hook() parses webhook and returns (Repo, Pipeline, error)
	//
	// Return Semantics:
	// - (repo, pipeline, nil): Execute pipeline for this event
	// - (repo, nil, nil): Valid webhook, no pipeline should run
	// - (nil, nil, types.ErrIgnoreEvent): Event ignored (logged)
	// - (nil, nil, error): Invalid webhook or parsing error
	//
	// Must verify webhook signature to prevent spoofing.
	// Should return types.ErrIgnoreEvent for non-pipeline events
	// (e.g. repository settings changed).
	Hook(ctx context.Context, r *http.Request) (*model.Repo, *model.Pipeline, error)


