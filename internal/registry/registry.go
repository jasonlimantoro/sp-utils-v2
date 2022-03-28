package registry

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gitlab"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/trello"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/repository"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/task"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createcard"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createlist"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/getweeklyupdates"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/listmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/reviewmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/syncrepo"
)

type Registry struct {
	CreateMergeRequestModule createmergerequest.Module
	ListMergeRequestModule   listmergerequest.Module
	ReviewMergeRequestModule reviewmergerequest.Module
	CreateCardModule         createcard.Module
	CreateListModule         createlist.Module
	SyncRepoModule           syncrepo.Module
	GetWeeklyUpdates         getweeklyupdates.Module
}

func InitRegistry() *Registry {
	reg := &Registry{}
	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	httpClient := &http.Client{}
	gitlabAccessor := gitlab.NewAccessor(httpClient)
	trelloAccessor := trello.NewAccessor(httpClient)

	repositoryDm := repository.NewManager(gitlabAccessor)
	taskDm := task.NewManager(trelloAccessor)

	reg.CreateMergeRequestModule = createmergerequest.NewModule(repositoryDm)
	reg.ListMergeRequestModule = listmergerequest.NewModule(repositoryDm)
	reg.ReviewMergeRequestModule = reviewmergerequest.NewModule(repositoryDm)
	reg.CreateCardModule = createcard.NewModule(taskDm, logrusLogger)
	reg.CreateListModule = createlist.NewModule(taskDm, logrusLogger)
	reg.SyncRepoModule = syncrepo.NewModule(logrusLogger)
	reg.GetWeeklyUpdates = getweeklyupdates.NewModule(taskDm, logrusLogger)

	return reg
}
