package xlog

import (
	"context"

	elog "github.com/labstack/gommon/log"

	"google.golang.org/appengine"
	alog "google.golang.org/appengine/log"
)

// Criticalf ...
func Criticalf(ctx context.Context, format string, args ...interface{}) {
	if !appengine.IsAppEngine() {
		elog.Panicf(format, args...)

		return
	}

	alog.Criticalf(ctx, format, args...)

	return
}

// Debugf ...
func Debugf(ctx context.Context, format string, args ...interface{}) {
	if !appengine.IsAppEngine() {
		elog.Debugf(format, args...)

		return
	}

	alog.Debugf(ctx, format, args...)

	return
}

// Errorf ...
func Errorf(ctx context.Context, format string, args ...interface{}) {
	if !appengine.IsAppEngine() {
		elog.Errorf(format, args...)

		return
	}

	alog.Errorf(ctx, format, args...)

	return
}

// Infof ...
func Infof(ctx context.Context, format string, args ...interface{}) {
	if !appengine.IsAppEngine() {
		elog.Infof(format, args...)

		return
	}

	alog.Infof(ctx, format, args...)

	return
}

// Warningf ...
func Warningf(ctx context.Context, format string, args ...interface{}) {
	if !appengine.IsAppEngine() {
		elog.Warnf(format, args...)

		return
	}

	alog.Warningf(ctx, format, args...)

	return
}
