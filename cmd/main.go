// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package main
package main

import (
	storageController "github.com/lijinfengnuc/prometheus-adapter/controller/storage"
	"github.com/lijinfengnuc/prometheus-adapter/flag"
	"github.com/lijinfengnuc/prometheus-adapter/router"
	storageService "github.com/lijinfengnuc/prometheus-adapter/service/storage"
	flagUtil "github.com/lijinfengnuc/prometheus-adapter/util/flag"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/sirupsen/logrus"
)

// main for build
func main() {
	//绑定命令行参数
	if err := flag.BindFlag(); err != nil {
		log.Logger.WithError(err).Error("bind flag error,exit")
		return
	}
	log.Logger.Info("bind flag success")
	//实例化adapter
	adapterName := *flagUtil.GetStringFlag(flag.AdapterName)
	adapterFilePath := *flagUtil.GetStringFlag(flag.AdapterFilePath)
	storage, err := storageService.GetStorage()

	if err != nil {
		log.Logger.WithError(err).WithFields(logrus.Fields{
			flag.AdapterName:     adapterName,
			flag.AdapterFilePath: adapterFilePath,
		}).Error("init storage error,exit")
		return
	}
	storageController.Storage = storage
	log.Logger.WithFields(logrus.Fields{
		flag.AdapterName:     adapterName,
		flag.AdapterFilePath: adapterFilePath,
	}).Info("init storage success")

	//绑定web服务
	router.BindAPI()
}
