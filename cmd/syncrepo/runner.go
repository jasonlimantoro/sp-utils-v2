package syncrepocmd

import (
	"context"

	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/syncrepo"
)

type runner struct {
	module syncrepo.Module
}

func NewRunner(module syncrepo.Module) *runner {
	return &runner{
		module: module,
	}
}

func (r runner) Run(ctx context.Context, flags map[string]string) error {
	return r.module.Do(ctx, (&syncrepo.Args{}).FromMap(flags))
}
