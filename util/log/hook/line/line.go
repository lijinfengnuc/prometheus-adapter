// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package line defines a hook for print line number
package line

import (
	"fmt"
	"github.com/lijinfengnuc/prometheus-adapter/util/log/hook"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
)

// LineHook inherits Hook and define some other fields
type LineHook struct {
	hook.Hook
	SourceFieldName string
	//This field will show func name
	ShowFunc      bool
	FuncFieldName string
	Skip          int
}

// NewDefault initialized a pointer of default LineHook
func NewDefault() *LineHook {
	return &LineHook{
		Hook: hook.Hook{
			Level: logrus.AllLevels,
		},
		SourceFieldName: "caller",
		ShowFunc:        false,
		FuncFieldName:   "func",
		Skip:            5,
	}
}

// New initialized a pointer of custom LineHook
func New(levels *[]logrus.Level, sourceFieldName string, showFunc bool, funcFieldName string, skip int) *LineHook {
	return &LineHook{
		Hook: hook.Hook{
			Level: *levels,
		},
		SourceFieldName: sourceFieldName,
		ShowFunc:        showFunc,
		FuncFieldName:   funcFieldName,
		Skip:            skip,
	}
}

// Levels implements interface Hook and returns levels
func (lineHook *LineHook) Levels() []logrus.Level {
	return lineHook.Level

}

// Fire implements interface Hook, adds line-num and func-name to the entry
func (lineHook *LineHook) Fire(entry *logrus.Entry) error {
	//调用runtime设置文件:行数、函数名称
	if pc, file, line, ok := runtime.Caller(lineHook.Skip); ok {
		//设置文件:行数
		entry.Data["caller"] = fmt.Sprintf("%s:%v", path.Base(file), line)
		//设置函数名称
		if lineHook.ShowFunc {
			funcName := runtime.FuncForPC(pc).Name()
			entry.Data["func"] = fmt.Sprintf("%s", path.Base(funcName))
		}
	} else {
		return errors.New("call runtime-caller error")
	}

	return nil
}
