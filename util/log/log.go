// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package log init a global var Logger for printing logs
package log

import (
	"github.com/lijinfengnuc/prometheus-adapter/util/log/hook/line"
	"github.com/sirupsen/logrus"
)

// Logger is a global var for printing logs in anywhere
var Logger *logrus.Entry

// init initialized a global var Logger
func init() {
	//设置时间格式
	timeFormatter := &logrus.TextFormatter{
		/*ForceColors:true,*/
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}
	logrus.SetFormatter(timeFormatter)

	//设置输出级别
	logrus.SetLevel(logrus.InfoLevel)

	//添加hook
	//添加line hook
	logrus.AddHook(line.NewDefault())

	//设置默认字段
	Logger = logrus.WithFields(logrus.Fields{})

	//打印成功信息
	Logger.Info("Logger is ready")
}
