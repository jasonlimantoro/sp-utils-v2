package syncrepocmd

import (
	"context"

	"git.garena.com/shopee/marketplace-payments/common/errlib"

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
	args, err := (&syncrepo.Args{}).FromMap(flags)
	if err != nil {
		return errlib.WrapFunc(err)
	}
	return r.module.Do(ctx, args)
}
