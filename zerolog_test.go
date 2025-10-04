// Copyright (c) 2025 Amari Robinson
// SPDX-License-Identifier: MIT

package fxeventzerolog

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"go.uber.org/fx/fxevent"
)

func newTestLogger() (*Logger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf)
	return New(&zl).(*Logger), buf
}

func TestLogger_DefaultLevels(t *testing.T) {
	logger, buf := newTestLogger()
	logger.log().Msg("info test")
	logger.err().Msg("error test")
	out := buf.String()
	if !strings.Contains(out, "info test") {
		t.Error("Expected info log message")
	}
	if !strings.Contains(out, "error test") {
		t.Error("Expected error log message")
	}
	if !strings.Contains(out, "\"level\":\"info\"") {
		t.Error("Expected info level for info log")
	}
	if !strings.Contains(out, "\"level\":\"error\"") {
		t.Error("Expected error level for error log")
	}
}

func TestLogger_CustomLevels(t *testing.T) {
	logger, buf := newTestLogger()
	logger.logLvl = zerolog.DebugLevel
	logger.errorLvl = zerolog.WarnLevel
	logger.log().Msg("debug test")
	logger.err().Msg("warn test")
	out := buf.String()
	if !strings.Contains(out, "debug test") {
		t.Error("Expected debug log message")
	}
	if !strings.Contains(out, "warn test") {
		t.Error("Expected warn log message")
	}
	if !strings.Contains(out, "\"level\":\"debug\"") {
		t.Error("Expected debug level for log")
	}
	if !strings.Contains(out, "\"level\":\"warn\"") {
		t.Error("Expected warn level for error log")
	}
}

func TestLogger_NilLoggerSafe(t *testing.T) {
	l := New(nil).(*Logger)
	// Should not panic
	l.LogEvent(&fxevent.OnStartExecuting{FunctionName: "f", CallerName: "c"})
}

func TestLogger_LogEvent_AllEvents(t *testing.T) {
	logger, buf := newTestLogger()
	events := []fxevent.Event{
		&fxevent.OnStartExecuting{FunctionName: "f1", CallerName: "c1"},
		&fxevent.OnStartExecuted{FunctionName: "f2", CallerName: "c2", Runtime: 123},
		&fxevent.OnStartExecuted{FunctionName: "f3", CallerName: "c3", Runtime: 456, Err: errors.New("fail")},
		&fxevent.OnStopExecuting{FunctionName: "f4", CallerName: "c4"},
		&fxevent.OnStopExecuted{FunctionName: "f5", CallerName: "c5", Runtime: 789},
		&fxevent.OnStopExecuted{FunctionName: "f6", CallerName: "c6", Runtime: 101, Err: errors.New("fail2")},
		&fxevent.Supplied{TypeName: "T", StackTrace: []string{"s1"}, ModuleTrace: []string{"m1"}},
		&fxevent.Supplied{TypeName: "T2", StackTrace: []string{"s2"}, ModuleTrace: []string{"m2"}, Err: errors.New("fail3")},
		&fxevent.Provided{ConstructorName: "ctor", OutputTypeNames: []string{"O1", "O2"}, StackTrace: []string{"s3"}, ModuleTrace: []string{"m3"}, Private: true},
		&fxevent.Provided{ConstructorName: "ctor2", OutputTypeNames: []string{"O3"}, StackTrace: []string{"s4"}, ModuleTrace: []string{"m4"}, Err: errors.New("fail4")},
		&fxevent.Run{Name: "run1", Kind: "kind1", Runtime: 1},
		&fxevent.Run{Name: "run2", Kind: "kind2", Runtime: 2, Err: errors.New("fail5")},
		&fxevent.Invoking{FunctionName: "fn1"},
		&fxevent.Invoked{FunctionName: "fn2", Trace: "trace1", Err: errors.New("fail6")},
		&fxevent.Stopping{Signal: os.Interrupt},
		&fxevent.Stopped{Err: errors.New("fail7")},
		&fxevent.RollingBack{StartErr: errors.New("fail8")},
		&fxevent.RolledBack{Err: errors.New("fail9")},
		&fxevent.Started{},
		&fxevent.Started{Err: errors.New("fail10")},
		&fxevent.LoggerInitialized{ConstructorName: "ctor3"},
		&fxevent.LoggerInitialized{Err: errors.New("fail11")},
	}
	for _, e := range events {
		logger.LogEvent(e)
	}
	out := buf.String()
	for _, want := range []string{
		"OnStart hook executing", "OnStart hook executed", "OnStart hook failed",
		"OnStop hook executing", "OnStop hook executed", "OnStop hook failed",
		"supplied", "provided", "error encountered while applying options",
		"run", "error returned", "invoking", "invoke failed",
		"received signal", "stop failed", "start failed", "started",
		"rolling back", "rollback failed", "initialized custom fxevent.Logger",
		"custom logger initialization failed",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Expected log output to contain %q", want)
		}
	}
}

func TestLogger_ModuleNameAndMaybeBool(t *testing.T) {
	buf := &bytes.Buffer{}
	zl := zerolog.New(buf)
	evt := zl.Info()
	evt = moduleName(evt, "mod1")
	evt = maybeBool(evt, "private", true)
	evt.Msg("test")
	out := buf.String()
	if !strings.Contains(out, "mod1") {
		t.Error("Expected module name in log output")
	}
	if !strings.Contains(out, "private") {
		t.Error("Expected private bool in log output")
	}
}
