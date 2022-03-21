package registry

import (
	"net/http"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gitlab"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/trello"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/repository"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/task"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createcard"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/listmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/reviewmergerequest"
)

type Registry struct {
	CreateMergeRequestModule createmergerequest.Module
	ListMergeRequestModule   listmergerequest.Module
	ReviewMergeRequestModule reviewmergerequest.Module
	CreateCardModule         createcard.Module
}

func InitRegistry() *Registry {
	reg := &Registry{}

	httpClient := &http.Client{}
	gitlabAccessor := gitlab.NewAccessor(httpClient)
	trelloAccessor := trello.NewAccessor(httpClient)

	repositoryDm := repository.NewManager(gitlabAccessor)
	taskDm := task.NewManager(trelloAccessor)

	reg.CreateMergeRequestModule = createmergerequest.NewModule(repositoryDm)
	reg.ListMergeRequestModule = listmergerequest.NewModule(repositoryDm)
	reg.ReviewMergeRequestModule = reviewmergerequest.NewModule(repositoryDm)
	reg.CreateCardModule = createcard.NewModule(taskDm)

	return reg
}
