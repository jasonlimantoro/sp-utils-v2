package createlistcmd

import (
	"context"

	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createlist"
)

type runner struct {
	module createlist.Module
}

func NewRunner(module createlist.Module) *runner {
	return &runner{module: module}
}

func (r runner) Run(ctx context.Context, flags map[string]string) error {
	return r.module.Do(ctx, (&createlist.Args{}).FromMap(flags))
}
