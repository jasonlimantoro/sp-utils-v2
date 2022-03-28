package logger

import "github.com/sirupsen/logrus"

type Fields = logrus.Fields

type Logger interface {
	logrus.FieldLogger
}
