package pr2doc

import (
    "context"
    "testing"

    "github.com/google/go-github/github"
    "github.com/pkg/errors"
)

type mockGithubService struct{}

// GetCommit is XXX
func (gs *mockGithubService) GetCommit(_ context.Context, _ string) (*github.RepositoryCommit, error) {
    return &github.RepositoryCommit{
        Commit: &github.Commit{
            Message: toPtr("Merge pull request #12345 from pr2doc/develop"),
        },
    }, nil
}

// GetPullRequest is XXX
func (gs *mockGithubService) GetPullRequestCommits(ctx context.Context, prNum int) ([]*github.RepositoryCommit, error) {
    if prNum != 12345 {
        return []*github.RepositoryCommit{}, errors.Errorf("error pull request #%d is not found", prNum)
    }
    return []*github.RepositoryCommit{
        {
            Commit: &github.Commit{
                Message: toPtr("Merge pull request #123 from pr2doc/develop"),
            },
        },
        {
            Commit: &github.Commit{
                Message: toPtr("Merge pull request #456 from pr2doc/develop"),
            },
        },
        {
            Commit: &github.Commit{
                Message: toPtr("Fix mobile nav menu"),
            },
        },
    }, nil
}

// GetPullRequest is XXX
func (gs *mockGithubService) GetPullRequest(_ context.Context, prNum int) (*github.PullRequest, error) {
    if prNum == 123 {
        return &github.PullRequest{
            Title: toPtr("Test title 1"),
            Body:  toPtr("This is test pull request.\n```share\nPlease shere this message 1\n```"),
        }, nil
    }
    if prNum == 345 {
        return &github.PullRequest{
            Title: toPtr("Test title 2"),
            Body:  toPtr("This is test pull request.\n```share\nPlease shere this message 2\n```"),
        }, nil
    }
    return nil, errors.Errorf("error pull request #%d is not found", prNum)
}

func toPtr(s string) *string {
    return &s
}

func Test_findDescription(t *testing.T) {
    type test struct {
        Title  string
        Input  string
        Output string
    }

    identifier := "test"
    tests := []test{
        {
            Title:  "success",
            Input:  "```test\ndescription\n```",
            Output: "description",
        },
        {
            Title:  "success (includes new line)",
            Input:  "```test\ndescription\ndescription\ndescription\n```",
            Output: "description\ndescription\ndescription",
        },
        {
            Title:  "error (mismatch identifier)",
            Input:  "```go\ntest\n```",
            Output: "",
        },
        {
            Title:  "error (missing new line)",
            Input:  "```test\ndescription```",
            Output: "",
        },
    }

    for _, test := range tests {
        t.Run(test.Title, func(t *testing.T) {
            p2d := &Pr2Doc{}
            if g, w := p2d.findDescription(test.Input, identifier), test.Output; g != w {
                t.Errorf("error find description %s, want %s", g, w)
            }
        })
    }
}

func Test_findPRNumber(t *testing.T) {
    type test struct {
        Title  string
        Input  string
        Output int
    }

    tests := []test{
        {
            Title:  "success",
            Input:  "Merge pull request #1234 from pr2doc/develop",
            Output: 1234,
        },
        {
            Title:  "success (too many number)",
            Input:  "Merge pull request #123 #567 #789 from pr2doc/develop",
            Output: 123,
        },
        {
            Title:  "error (missing PR number requires # )",
            Input:  "Merge pull request 12345 from pr2doc/develop",
            Output: 0,
        },
        {
            Title:  "error (no spacing)",
            Input:  "Merge pull request#987 from pr2doc/develop",
            Output: 0,
        },
        {
            Title:  "error (too large number)",
            Input:  "Merge pull request #456789 from pr2doc/develop",
            Output: 0,
        },
    }

    for _, test := range tests {
        t.Run(test.Title, func(t *testing.T) {
            p2d := &Pr2Doc{}
            prNum, err := p2d.findPRNumber(test.Input)
            if err != nil {
                t.Fatal("error findPRNumber", err.Error())
            }
            if g, w := prNum, test.Output; g != w {
                t.Errorf("error find PR number %d, want %d", g, w)
            }
        })
    }
}
