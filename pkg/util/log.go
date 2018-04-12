package util

import (
	"github.com/golang/glog"
)

// Debug prints message with severity 3.
func Debug(f string, args ...interface{}) {
	if glog.V(3) {
		glog.Infof(f, args...)
	}
}

// Info prints message with severity 2.
func Info(f string, args ...interface{}) {
	if glog.V(2) {
		glog.Infof(f, args...)
	}
}

// Error prints error message.
func Error(f string, args ...interface{}) {
	glog.Errorf(f, args...)
}

// Exit prints message and exit program with code 1.
func Exit(f string, args ...interface{}) {
	glog.Exitf(f, args...)
}

// Flush flushes glog stream.
func Flush() {
	glog.Flush()
}
