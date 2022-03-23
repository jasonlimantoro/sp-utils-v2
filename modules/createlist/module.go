package createlist

import (
	"context"
	"strconv"
	"time"

	"git.garena.com/shopee/marketplace-payments/common/errlib"
	"github.com/sirupsen/logrus"

	"git.garena.com/jason.limantoro/shopee-utils-v2/internal/manager/task"
)

const (
	OperationTypeToday = iota
	OperationTypeNextWorkingDay
)

type Module interface {
	Do(ctx context.Context, Args *Args) error
}

type module struct {
	manager task.Manager
	logger  *logrus.Logger
}

func NewModule(manager task.Manager, logger *logrus.Logger) *module {
	return &module{manager: manager, logger: logger}
}

func (m module) Do(ctx context.Context, args *Args) error {
	var name string
	today := getTodayString()
	switch args.OperationType {
	case OperationTypeToday:
		name = today
	case OperationTypeNextWorkingDay:
		name = getNextWorkingDayString()
	default:
		name = today
	}

	list, err := m.manager.CreateList(ctx, name)

	if err != nil {
		return errlib.WrapFunc(err)
	}

	m.logger.WithFields(logrus.Fields{
		"name": list.Name,
		"pos":  list.Pos,
	}).Info("created")

	return nil
}

func getNextWorkingDayString() string {
	currentTime := time.Now()
	var deltaDays int
	if currentTime.Weekday() == time.Friday {
		deltaDays = 3
	} else {
		deltaDays = 1
	}

	return currentTime.AddDate(0, 0, deltaDays).Format("02-Jan-2006")
}

func getTodayString() string {
	return time.Now().Format("02-Jan-2006")
}

type Args struct {
	OperationType int
}

func (a *Args) FromMap(flags map[string]string) *Args {
	operationTypeInt, _ := strconv.Atoi(flags["operation-type"])
	a.OperationType = operationTypeInt

	return a
}
