package registry

import (
	"net/http"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gitlab"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/repository"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createmergerequest"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/listmergerequest"
)

type Registry struct {
	CreateMergeRequestModule createmergerequest.Module
	ListMergeRequestModule   listmergerequest.Module
}

func InitRegistry() *Registry {
	reg := &Registry{}

	httpClient := &http.Client{}
	gitlabAccessor := gitlab.NewAccessor(httpClient)

	repositoryDm := repository.NewManager(gitlabAccessor)

	reg.CreateMergeRequestModule = createmergerequest.NewModule(repositoryDm)
	reg.ListMergeRequestModule = listmergerequest.NewModule(repositoryDm)

	return reg
}
