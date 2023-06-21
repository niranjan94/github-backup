package internal

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/v53/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
	"os"
	"os/signal"
	"syscall"
)

type Backup struct {
	Config *Config
}

func (b *Backup) getRepositoryUrl(userLogin string, repository *github.Repository) string {
	return fmt.Sprintf("https://%s:%s@github.com/%s/%s.git", userLogin, b.Config.Token, *repository.Owner.Login, *repository.Name)
}

func (b *Backup) downloadRepository(ctx context.Context, userLogin string, repository *github.Repository) error {
	repositoryPathComponent := path.Join(*repository.Owner.Login, *repository.Name)
	repositoryPath := path.Join(b.Config.Directory, repositoryPathComponent)
	if err := os.MkdirAll(repositoryPath, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %v", repositoryPath, err)
	}

	logrus.Infof("Backing up %s to %s", repositoryPathComponent, repositoryPath)
	localRepository, err := git.PlainOpenWithOptions(repositoryPath, &git.PlainOpenOptions{
		DetectDotGit: false,
	})
	if err != nil {
		logrus.Infoln("Cloning repository", repositoryPathComponent)
		_, err := git.PlainCloneContext(ctx, repositoryPath, true, &git.CloneOptions{
			URL:      b.getRepositoryUrl(userLogin, repository),
			Progress: nil,
			Mirror:   true,
		})
		if err != nil {
			return fmt.Errorf("error cloning repository %s: %v", *repository.Name, err)
		}
	} else {
		logrus.Infoln("Updating repository", repositoryPathComponent)
		if err := localRepository.FetchContext(ctx, &git.FetchOptions{
			Force: true,
			Tags:  git.AllTags,
			RefSpecs: []config.RefSpec{
				"refs/heads/*:refs/heads/*",
			},
		}); err != nil {
			if err != git.NoErrAlreadyUpToDate {
				return fmt.Errorf("error fetching repository %s: %v", *repository.Name, err)
			}
		}
		logrus.Infoln("Pruning repository", repositoryPathComponent)
		if err := localRepository.Prune(git.PruneOptions{}); err != nil {
			return fmt.Errorf("error pruning repository %s: %v", *repository.Name, err)
		}
	}
	return nil
}

func (b *Backup) perform(ctx context.Context) error {
	client := github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: b.Config.Token},
	)))

	req, err := client.NewRequest("GET", "/user", nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	var user github.User
	_, err = client.Do(ctx, req, &user)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	userLogin := *user.Login

	logrus.Infof("Using GitHub user: %s", userLogin)
	logrus.Infof("Loading repositories...")

	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 1000},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.List(ctx, "", opt)
		if err != nil {
			return fmt.Errorf("error listing repositories: %v", err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	logrus.Infof("Found %d repositories", len(allRepos))
	logrus.Infof("Filtering repositories by owner: %s", b.Config.Owners)

	var filteredRepos []*github.Repository
	for idx, repo := range allRepos {
		if slices.Contains(b.Config.Owners, strings.ToLower(*repo.Owner.Login)) {
			filteredRepos = append(filteredRepos, allRepos[idx])
		}
	}

	logrus.Infof("Found %d repositories after filtering", len(filteredRepos))
	for _, repo := range filteredRepos {
		if ctx.Err() != nil {
			return nil
		}
		if err := b.downloadRepository(ctx, userLogin, repo); err != nil {
			logrus.Errorf("error downloading repository: %v", err)
		}
	}
	return nil
}

func (b *Backup) Start() {
	// Create a new context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Cancel ctx as soon as handleSearch is done

	// Handle SIGINT and SIGTERM
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case <-ch:
			logrus.Infoln("Received shutdown signal")
			cancel()
		}
	}()

	for {
		// Check if context is cancelled, if so, return
		if ctx.Err() != nil {
			logrus.Infoln("Context cancelled")
			return
		}

		// Perform the backup
		if err := b.perform(ctx); err != nil {
			logrus.Errorf("error performing backup: %v", err)
		}

		// Sleep for the period
		logrus.Infof("Sleeping for %d seconds", b.Config.PeriodSeconds)
		DelayWithContext(ctx, time.Duration(b.Config.PeriodSeconds)*time.Second)
	}

}

func NewBackupSchedule(config *Config) *Backup {
	return &Backup{config}
}
