package logging

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/mdev5000/qlik_message/apperrors"
	"github.com/sirupsen/logrus"
	"os"
)

const LogInfo = true
const LogWarn = true

type Logger struct {
	logrus.Logger
}

func (l *Logger) LogFailedToEncode(op string, origErr error, jsonErr error, stackErr error) {
	l.WithFields(Fields{
		"originalErr:": fmt.Sprintf("%+v", origErr),
		"jsonErr":      fmt.Sprintf("%+v", jsonErr),
		"stack":        fmt.Sprintf("%+v", stackErr),
	}).Errorf("failed to encode error response (op: %s)", op)
}

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
