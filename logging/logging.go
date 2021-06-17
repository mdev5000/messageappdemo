// Package logging contains core logging logic for the application. Currently mostly just a wrapper around the logrus
// package.
package logging

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/mdev5000/messageappdemo/apperrors"
	"github.com/sirupsen/logrus"
)

const LogInfo = true
const LogWarn = true

type Logger struct {
	logrus.Logger
}

// LogFailedToEncode indicates if an application specific json struct failed to encode. Ideally you should never see
// this logged.
func (l *Logger) LogFailedToEncode(op string, origErr error, jsonErr error, stackErr error) {
	l.WithFields(Fields{
		"originalErr:": fmt.Sprintf("%+v", origErr),
		"jsonErr":      fmt.Sprintf("%+v", jsonErr),
		"stack":        fmt.Sprintf("%+v", stackErr),
	}).Errorf("failed to encode error response (op: %s)", op)
}

// LogError an error from within the application. If the error is of type apperrors.Error then the information inside
// is specially encoded into a log entry.
func (l *Logger) LogError(err error) {
	switch e := err.(type) {
	case *apperrors.Error:
		l.WithFields(Fields{
			"err":   fmt.Sprintf("%+v", e.Err),
			"stack": fmt.Sprintf("%+v", e.Stack),
		}).Errorf("%s error (op: %s)", e.EType, e.Op)
	default:
		l.WithFields(Fields{
			"err": fmt.Sprintf("%+v", err),
		}).Errorf("internal error")
	}
}

type Fields = logrus.Fields

// Dump creates a string representation of a value. Using the spew library to provide more detail.
func Dump(value ...interface{}) string {
	return spew.Sdump(value...)
}

func NoLog() *Logger {
	l := Logger{
		Logger: *logrus.New(),
	}
	l.Logger.SetOutput(ioutil.Discard)
	return &l
}

func New() *Logger {
	l := Logger{
		Logger: *logrus.StandardLogger(),
	}
	l.Logger.SetLevel(logrus.WarnLevel)
	l.Logger.SetOutput(os.Stdout)
	l.Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return &l
}
