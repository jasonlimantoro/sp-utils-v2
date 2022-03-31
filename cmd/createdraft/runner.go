package createdraftcmd

import (
	"context"

	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createdraft"
)

type runner struct {
	module createdraft.Module
}

func (r runner) Run(ctx context.Context, flags map[string]string) error {
	return r.module.Do(ctx, (&createdraft.Args{}).FromMap(flags))
}

func NewRunner(module createdraft.Module) *runner {
	return &runner{module: module}
}
