package syncrepo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"git.garena.com/shopee/marketplace-payments/common/errlib"
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
)

var (
	stagingStatusAllowed = map[git.StatusCode]bool{
		git.Unmodified:         true,
		git.Untracked:          true,
		git.Added:              true,
		git.Modified:           false,
		git.Renamed:            false,
		git.Deleted:            false,
		git.UpdatedButUnmerged: false,
	}

	worktreeStatusAllowed = map[git.StatusCode]bool{
		git.Unmodified:         true,
		git.Untracked:          true,
		git.Added:              true,
		git.Modified:           false,
		git.Renamed:            false,
		git.Deleted:            false,
		git.UpdatedButUnmerged: false,
	}
)

const (
	RemoteOrigin = "origin"
)

type Module interface {
	Do(ctx context.Context, args *Args) error
}

type module struct {
	logger logrus.FieldLogger
}

func NewModule(logger logrus.FieldLogger) *module {
	return &module{
		logger: logger,
	}
}

func (m module) Do(ctx context.Context, args *Args) error {
	start := time.Now()
	argBytes, _ := json.Marshal(args)
	cmdLogger := m.logger.WithFields(logrus.Fields{
		"args": string(argBytes),
	})
	cmdLogger.Info("start")

	defer func() {
		elapsed := time.Since(start)
		cmdLogger.WithFields(logrus.Fields{
			"elapsed": elapsed.Seconds(),
		}).Info("ended")
	}()

	wg := &sync.WaitGroup{}
	wg.Add(len(args.Repositories))
	for _, repository := range args.Repositories {
		go func(repository Repository) {
			defer wg.Done()
			if err := m.process(ctx, repository.Path, repository.TargetBranches); err != nil {
				m.logger.WithError(errlib.WithFields(err, errlib.Fields{
					"repository": repository.Path,
				})).Error("error processing")
			}
		}(repository)
	}
	wg.Wait()

	return nil
}

func (m module) process(ctx context.Context, repositoryPath string, branches []Branch) error {
	repositoryLogger := m.logger.WithFields(logrus.Fields{
		"repository": repositoryPath,
	})

	r, err := git.PlainOpen(repositoryPath)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		repositoryLogger.WithError(err).Warn("Skipping directory")
		return nil
	}
	if err != nil {
		return errlib.WrapFunc(err)
	}

	h, err := r.Head()
	if err != nil {
		return errlib.WrapFunc(err)
	}
	initialBranch := h.Name().Short()

	w, err := r.Worktree()
	if err != nil {
		return errlib.WrapFunc(err)
	}

	status, err := w.Status()
	if err != nil {
		return errlib.WrapFunc(err)
	}

	shouldSync := true
	for _, s := range status {
		isFileStatusAllowed := stagingStatusAllowed[s.Staging] && worktreeStatusAllowed[s.Worktree]
		if !isFileStatusAllowed {
			shouldSync = false
		}
	}

	repositoryLogger.Infof("shouldSync=%v", shouldSync)

	if shouldSync {
		for _, branch := range branches {
			branchLogger := repositoryLogger.WithField("branch", branch)

			branchLogger.WithFields(logrus.Fields{
				"command": "git checkout",
			}).Info("start")
			// Note: the checkout implementation of go-git is different from that of git.
			// Open issue: https://github.com/src-d/go-git/issues/1026
			out, err := exec.Command("git", "-C", repositoryPath, "checkout", branch.Local).CombinedOutput()
			if err != nil {
				return errlib.WrapFunc(errlib.WithFields(err, errlib.Fields{
					"command": "git checkout",
					"branch":  branch,
					"out":     string(out),
				}))
			}
			branchLogger.WithFields(logrus.Fields{
				"command": "git checkout",
			}).Info("success")

			branchLogger.WithFields(logrus.Fields{
				"command": fmt.Sprintf("git pull %s", RemoteOrigin),
			}).Info("start")
			start := time.Now()
			// Note: the pull implementation of go-git is super slow
			if out, err = exec.Command("git", "-C", repositoryPath, "pull", RemoteOrigin, branch.Remote).CombinedOutput(); err != nil {
				return errlib.WrapFunc(errlib.WithFields(err, errlib.Fields{
					"branch":  branch,
					"command": fmt.Sprintf("git pull %s", RemoteOrigin),
					"out":     string(out),
				}))
			}
			elapsed := time.Since(start)
			branchLogger.WithFields(logrus.Fields{
				"elapsed": elapsed.Seconds(),
				"command": fmt.Sprintf("git pull %s", RemoteOrigin),
				"result":  string(out),
			}).Info("success")
		}

		repositoryLogger.WithFields(logrus.Fields{
			"branch":  initialBranch,
			"command": "git checkout <initial_branch>",
		}).Info("start")

		out, err := exec.Command("git", "-C", repositoryPath, "checkout", initialBranch).CombinedOutput()
		if err != nil {
			return errlib.WrapFunc(errlib.WithFields(err, errlib.Fields{
				"command": "git checkout",
				"branch":  initialBranch,
				"out":     string(out),
			}))
		}

		repositoryLogger.WithFields(logrus.Fields{
			"branch":  initialBranch,
			"command": "git checkout <initial_branch>",
		}).Info("success")
	}

	return nil
}

type Args struct {
	Repositories []Repository
}

type Repository struct {
	Path           string
	TargetBranches []Branch
}

type Branch struct {
	Local  string
	Remote string
}

func (a *Args) FromMap(flags map[string]string) (*Args, error) {
	repoFile := flags["repo-file"]
	if repoFile != "" {
		path, err := filepath.Abs(repoFile)
		if err != nil {
			return nil, errlib.WrapFunc(err)
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errlib.WrapFunc(err)
		}

		rows := strings.Split(string(content), "\n")
		for _, row := range rows {
			cols := strings.Split(row, ",")
			path := cols[0]
			branchPairs := strings.Split(cols[1], ";")

			branches := []Branch{}
			for _, branchPair := range branchPairs {
				branchLocalRemote := strings.Split(branchPair, ":")
				branches = append(branches, Branch{
					Local:  branchLocalRemote[0],
					Remote: branchLocalRemote[1],
				})
			}

			a.Repositories = append(a.Repositories, Repository{
				Path:           path,
				TargetBranches: branches,
			})
		}
	}

	return a, nil
}
