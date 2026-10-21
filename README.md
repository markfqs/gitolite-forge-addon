Just compile/test:

go mod init generic-forge && go mod tidy && go build -o generic ./generic


TO IMPLEMENT
+	Name() string
+	URL() string
!	Login(ctx context.Context, r *types.OAuthRequest) (*model.User, string, error)
-	Teams(ctx context.Context, u *model.User, p *model.ListOptions) ([]*model.Team, error)
	Hook(ctx context.Context, r *http.Request) (*model.Repo, *model.Pipeline, error)
	Repo(ctx context.Context, u *model.User, remoteID model.ForgeRemoteID, owner, name string) (*model.Repo, error)
	Repos(ctx context.Context, u *model.User, p *model.ListOptions) ([]*model.Repo, error)
	File(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, fileName string) ([]byte, error)
-	Dir(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, dirName string) ([]*types.FileMeta, error)
-	Status(ctx context.Context, u *model.User, r *model.Repo, b *model.Pipeline, p *model.Workflow) error
-	Netrc(u *model.User, r *model.Repo) (*model.Netrc, error)
-	Activate(ctx context.Context, u *model.User, r *model.Repo, link string) error
-	Deactivate(ctx context.Context, u *model.User, r *model.Repo, link string) error
-	Branches(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]string, error)
-	BranchHead(ctx context.Context, u *model.User, r *model.Repo, branch string) (*model.Commit, error)
-	PullRequests(ctx context.Context, u *model.User, r *model.Repo, p *model.ListOptions) ([]*model.PullRequest, error)
-	OrgMembership(ctx context.Context, u *model.User, org string) (*model.OrgPerm, error)
-	Org(ctx context.Context, u *model.User, org string) (*model.Org, error)


As we need to manage different Auth URLs from the repos, we must get/set them from enviroment and add to config/forge struct/class

PROMPT
> I'm developing a WoodPeckerCI forge-addon that will use a generic OAuth2 provider (so all needed parameters like URL and the like need to be configurable). Git repositories are accessed via SSH with a dedicated, fixed SSH-key no matter the user and the pipelines will be triggered by post-receive hooks on the git repositories. Tell me the minimal forge methods interface I need to implement.

Any authenticated user should have access to the pipelines (so I just need to authenticate, no )


Required methods:
Login -> ?
Repo -> Return a forge repository URL
Repos -> List of repositories
File
Hook

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



Repository: woodpecker-generic-forge
Forge identifier: generic
Go module: codeberg.org/yourusername/woodpecker-generic-forge (if hosting on Codeberg)
Executable: generic

