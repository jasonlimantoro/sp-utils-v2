package registry

import (
	"net/http"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/accessor/gitlab"
	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/repository"
	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createmergerequest"
)

type Registry struct {
	CreateMergeRequestModule createmergerequest.Module
}

func InitRegistry() *Registry {
	reg := &Registry{}

	httpClient := &http.Client{}
	reg.CreateMergeRequestModule = createmergerequest.NewModule(
		repository.NewManager(gitlab.NewAccessor(httpClient)),
	)

	return reg
}
