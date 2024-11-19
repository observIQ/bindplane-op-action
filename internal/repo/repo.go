package repo

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// CloneRepo clones a repository from the provided URL and branch. If the cloneURL
// is empty, the function will attempt to clone the repository using a GitHub
// URL assembled from the GITHUB_ACTOR, GITHUB_REPOSITORY environment variables,
// and the provided token.
func CloneRepo(cloneURL, branch, token string) (*git.Repository, error) {
	// TODO(jsirianni): githubURL should be assembled outside of this package
	// and passed in. We could remove the need for token, actor, and repo.

	githubActor := os.Getenv("GITHUB_ACTOR")
	githubRepo := os.Getenv("GITHUB_REPOSITORY")
	if cloneURL == "" {
		cloneURL = fmt.Sprintf(
			"https://%s:%s@github.com/%s.git",
			githubActor,
			token,
			githubRepo,
		)
	}

	// TODO(jsirianni): This context sets the clone timeout. This should be
	// a configurable option.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	dir := fmt.Sprintf("%s/%s-%d", os.TempDir(), branch, time.Now().Unix())

	repo, err := git.PlainCloneContext(ctx, dir, false, &git.CloneOptions{
		URL:           cloneURL,
		Progress:      os.Stdout,
		SingleBranch:  true,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
	})
	if err != nil {
		return nil, fmt.Errorf("clone repository: %w", err)
	}

	return repo, nil
}
