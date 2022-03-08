package listmergerequestcmd

import (
	"context"

	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/listmergerequest"
)

type runner struct {
	module listmergerequest.Module
}

func NewRunner(module listmergerequest.Module) *runner {
	return &runner{
		module: module,
	}
}

func (r runner) Run(ctx context.Context, flags map[string]string) error {
	return r.module.Do(ctx, (&listmergerequest.Args{}).FromMap(flags))
}
