package createmergerequestcmd

import (
	"context"

	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createmergerequest"
)

type runner struct {
	module createmergerequest.Module
}

func NewRunner(module createmergerequest.Module) *runner {
	return &runner{
		module: module,
	}
}

func (r runner) Run(ctx context.Context, flags map[string]string) error {
	return r.module.Do(ctx, (&createmergerequest.Args{}).FromMap(flags))
}
