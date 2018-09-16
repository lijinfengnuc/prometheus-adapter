// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package router binds APIs and listens specified port
package router

import (
	"github.com/gin-gonic/gin"
	"github.com/lijinfengnuc/prometheus-adapter/controller/health"
	"github.com/lijinfengnuc/prometheus-adapter/controller/storage"
	"github.com/lijinfengnuc/prometheus-adapter/flag"
	flagUtil "github.com/lijinfengnuc/prometheus-adapter/util/flag"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
)

// BindAPI binds APIs and listens specified port
func BindAPI() {
	//实例化router
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)
	router.Use(gin.Recovery())

	//绑定API
	v1 := router.Group("/v1")
	{
		//绑定health接口
		v1.GET("/health", health.Health)
		//绑定存储、读取指标接口
		v1.POST("/read", storage.Read)
		v1.POST("/write", storage.Write)
	}

	//指定端口启动web服务
	port := ":" + *flagUtil.GetStringFlag(flag.WebListenPort)
	log.Logger.Info("router start on port" + port)
	router.Run(port)
}
