package syncrepo

import (
	"context"
	"errors"
	"strings"
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
		"directory": args.Directory,
	})
	cmdLogger.Info("start")

	defer func() {
		elapsed := time.Since(start)
		cmdLogger.WithFields(logrus.Fields{
			"elapsed": elapsed.Seconds(),
		}).Info("ended")
	}()

	r, err := git.PlainOpen(args.Directory)
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
	for _, fs := range status {
		isFileStatusAllowed := stagingStatusAllowed[fs.Staging] && worktreeStatusAllowed[fs.Worktree]
		if !isFileStatusAllowed {
			shouldSync = false
		}
	}

	cmdLogger.Infof("shouldSync=%v", shouldSync)

	if shouldSync {
		for _, branch := range args.Branches {
			branchLogger := cmdLogger.WithField("branch", branch)

			branchLogger.Infof("[start] git checkout")
			if err := w.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(branch)}); err != nil {
				return errlib.WrapFunc(err)
			}
			branchLogger.Info("[success] git checkout")

			branchLogger.Info("[start] git pull origin")
			start := time.Now()
			if err := w.PullContext(ctx, &git.PullOptions{
				RemoteName:    "origin",
				ReferenceName: plumbing.NewBranchReferenceName(branch),
			}); err != nil {
				if errors.Is(err, git.NoErrAlreadyUpToDate) {
					branchLogger.Info("git pull origin: already up to date")
				} else {
					return errlib.WrapFunc(err)
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
	Directory string
	Branches  []string
}

func (a *Args) FromMap(flags map[string]string) *Args {
	a.Directory = flags["directory"]
	a.Branches = strings.Split(flags["branch"], ",")

	return a
}
