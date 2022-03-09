package reviewmergerequestcmd

import (
	"context"

	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/reviewmergerequest"
)

type runner struct {
	module reviewmergerequest.Module
}

func NewRunner(module reviewmergerequest.Module) *runner {
	return &runner{
		module: module,
	}
}

func (r runner) Run(ctx context.Context, flags map[string]string) error {
	return r.module.Do(ctx, (&reviewmergerequest.Args{}).FromMap(flags))
}
