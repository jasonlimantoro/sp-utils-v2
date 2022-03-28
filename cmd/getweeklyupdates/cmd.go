package getweeklyupdatescmd

import (
	"context"

	"git.garena.com/jason.limantoro/shopee-utils-v2/modules/getweeklyupdates"
)

type runner struct {
	module getweeklyupdates.Module
}

func NewRunner(module getweeklyupdates.Module) *runner {
	return &runner{module: module}
}

func (r runner) Run(ctx context.Context, flags map[string]string) error {
	return r.module.Do(ctx, (&getweeklyupdates.Args{}).FromMap(flags))
}
