package syncrepo

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"git.garena.com/shopee/marketplace-payments/common/errlib"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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
	errChan := make(chan error)
	doneChan := make(chan bool)

	wg.Add(len(allRepositories))

	for _, repository := range allRepositories {
		go func(repository string) {
			defer wg.Done()
			if err := m.process(ctx, repository, args.Branches); err != nil {
				errChan <- errlib.WithFields(err, errlib.Fields{
					"repository": repository,
				})
			}
		}(repository)
	}

	go func() {
		wg.Wait()
		close(doneChan)
	}()

	for i := 0; i < len(allRepositories); i++ {
		select {
		case err := <-errChan:
			m.logger.WithError(err).Error("error processing")
			break
		case <-doneChan:
			break
		}
	}

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

func (m module) process(ctx context.Context, directoryPath string, branches []string) error {
	directoryLogger := m.logger.WithFields(logrus.Fields{
		"directory": directoryPath,
	})

	r, err := git.PlainOpen(directoryPath)
	if errors.Is(err, git.ErrRepositoryNotExists) {
		directoryLogger.WithError(err).Warn("Skipping directory")
		return nil
	}

	if err != nil {
		return errlib.WrapFunc(err)
	}

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

	directoryLogger.Infof("shouldSync=%v", shouldSync)

	if shouldSync {
		for _, branch := range branches {
			branchLogger := directoryLogger.WithField("branch", branch)

			branchLogger.Info("[start] git checkout")
			if err := w.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(branch)}); err != nil {
				if errors.Is(err, plumbing.ErrReferenceNotFound) {
					branchLogger.Warn("[warn] git checkout: branch not found")
					continue
				}

				return errlib.WrapFunc(errlib.WithFields(err, errlib.Fields{
					"branch": branch,
				}))
			}
			branchLogger.Info("[success] git checkout")

			branchLogger.Info("[start] git pull origin")
			start := time.Now()
			if err := w.PullContext(ctx, &git.PullOptions{
				RemoteName:    "origin",
				ReferenceName: plumbing.NewBranchReferenceName(branch),
			}); err != nil {
				if errors.Is(err, git.NoErrAlreadyUpToDate) {
					branchLogger.Warn("[warn] git pull origin: already up to date")
				} else {
					return errlib.WrapFunc(errlib.WithFields(err, errlib.Fields{
						"branch": branch,
					}))
				}
			}

			elapsed := time.Since(start)

			branchLogger.WithFields(logrus.Fields{
				"elapsed": elapsed.Seconds(),
			}).Info("[success] git pull origin")
		}
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
