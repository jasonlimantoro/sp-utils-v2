package syncrepo

import (
	"context"
	"errors"
	"io/fs"
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

type Module interface {
	Do(ctx context.Context, args *Args) error
}

type module struct {
	logger *logrus.Logger
}

func NewModule(logger *logrus.Logger) *module {
	return &module{
		logger: logger,
	}
}

func (m module) Do(ctx context.Context, args *Args) error {
	start := time.Now()
	cmdLogger := m.logger.WithFields(logrus.Fields{
		"root":     args.Root,
		"branches": args.Branches,
	})
	cmdLogger.Info("start")

	defer func() {
		elapsed := time.Since(start)
		cmdLogger.WithFields(logrus.Fields{
			"elapsed": elapsed.Seconds(),
		}).Info("ended")
	}()

	allRepositories := []string{}
	for _, root := range args.Root {
		repositories, err := findRepositories(root)
		if err != nil {
			return errlib.WrapFunc(err)
		}
		allRepositories = append(allRepositories, repositories...)
	}

	m.logger.Infof("repositories: %+v", allRepositories)

	wg := &sync.WaitGroup{}
	wg.Add(len(allRepositories))
	for _, repository := range allRepositories {
		go func(repository string) {
			defer wg.Done()
			if err := m.process(ctx, repository, args.Branches); err != nil {
				m.logger.WithError(errlib.WithFields(err, errlib.Fields{
					"repository": repository,
				})).Error("error processing")
			}
		}(repository)
	}
	wg.Wait()

	return nil
}

func findRepositories(rootDirectory string) ([]string, error) {
	res := []string{}
	err := filepath.WalkDir(rootDirectory, func(path string, d fs.DirEntry, err error) error {
		_, plainOpenErr := git.PlainOpen(path)
		if plainOpenErr == nil {
			res = append(res, path)
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, errlib.WrapFunc(err)
	}

	return res, nil
}

func (m module) process(ctx context.Context, repositoryPath string, branches []string) error {
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
			out, err := exec.Command("git", "-C", repositoryPath, "checkout", branch).CombinedOutput()
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
				"command": "git pull",
			}).Info("start")
			start := time.Now()
			// Note: the pull implementation of go-git is super slow
			if out, err = exec.Command("git", "-C", repositoryPath, "pull").CombinedOutput(); err != nil {
				return errlib.WrapFunc(errlib.WithFields(err, errlib.Fields{
					"branch":  branch,
					"command": "git pull",
					"out":     string(out),
				}))
			}
			elapsed := time.Since(start)
			branchLogger.WithFields(logrus.Fields{
				"elapsed": elapsed.Seconds(),
				"command": "git pull",
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
	Root     []string
	Branches []string
}

func (a *Args) FromMap(flags map[string]string) *Args {
	a.Root = strings.Split(flags["root"], ",")
	a.Branches = strings.Split(flags["branch"], ",")

	return a
}
