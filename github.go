package pr2doc

import (
	"context"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

type githubService struct {
	client *github.Client
	owner  string
	repo   string
}

func NewGithubService(ctx context.Context, owner, repo, token string) GithubService {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &githubService{
		client: github.NewClient(tc),
		owner:  owner,
		repo:   repo,
	}
}

// GithubService はGithubへの操作を行うinterface
type GithubService interface {
	GetCommit(ctx context.Context, sha string) (*github.RepositoryCommit, error)
	GetPullRequest(ctx context.Context, prNum int) (*github.PullRequest, error)
	GetPullRequestCommits(ctx context.Context, prNum int) ([]*github.RepositoryCommit, error)
}

// GetCommit is XXX
func (gs *githubService) GetCommit(ctx context.Context, sha string) (*github.RepositoryCommit, error) {
	cmt, res, err := gs.client.Repositories.GetCommit(ctx, gs.owner, gs.repo, sha)
	if err != nil {
		return nil, errors.Wrap(err, "Repositories.GetCommit failed")
	}
	defer res.Body.Close()
	return cmt, nil
}

// GetPullRequest is XXX
func (gs *githubService) GetPullRequest(ctx context.Context, prNum int) (*github.PullRequest, error) {
	pr, res, err := gs.client.PullRequests.Get(ctx, gs.owner, gs.repo, prNum)
	if err != nil {
		return nil, errors.Wrap(err, "PullRequest.Get failed")
	}
	defer res.Body.Close()
	return pr, nil
}

// GetPullRequest is XXX
func (gs *githubService) GetPullRequestCommits(ctx context.Context, prNum int) ([]*github.RepositoryCommit, error) {
	commits, res, err := gs.client.PullRequests.ListCommits(ctx, gs.owner, gs.repo, prNum, nil)
	if err != nil {
		return nil, errors.Wrap(err, "PullRequest.ListCommits failed")
	}
	defer res.Body.Close()
	return commits, nil
}
