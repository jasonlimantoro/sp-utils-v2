package createcardcmd

import (
	"context"

	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/createcard"
)

type runner struct {
	module createcard.Module
}

func NewRunner(module createcard.Module) *runner {
	return &runner{module: module}
}

func (r runner) Run(ctx context.Context, flags map[string]string) error {
	return r.module.Do(ctx, (&createcard.Args{}).FromMap(flags))
}
