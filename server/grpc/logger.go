package grpc

import (
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func newTmLogger(logger tmlog.Logger, msg string) *tmLogger {
	return &tmLogger{
		Logger: logger,
		msg:    msg,
	}
}

type tmLogger struct {
	tmlog.Logger
	msg string
}

func (l tmLogger) Log(keyvals ...interface{}) error {
	l.Info(l.msg, keyvals...)
	return nil
}
