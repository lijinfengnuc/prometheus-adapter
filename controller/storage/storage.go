// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package storage defines Read/Write controller for prometheus
package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/lijinfengnuc/prometheus-adapter/service/storage"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/lijinfengnuc/prometheus-adapter/util/prometheus"
	"github.com/prometheus/prometheus/prompb"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// -- Some constants
const (
	ReadPath  = "/read"
	WritePath = "/write"
	Path      = "path"
)

// Storage is a var of interface Storage
var Storage storage.Storage

// Read is a controller to query metrics from storage
func Read(ctx *gin.Context) {
	begin := time.Now()
	//打印日志
	log.Logger.WithFields(logrus.Fields{
		Path: ReadPath,
	}).Info("receive request from prometheus")
	//解码request
	request := &prompb.ReadRequest{}
	if err := prometheus.Unmarshal(request, ctx.Request); err != nil {
		log.Logger.WithError(err).WithFields(logrus.Fields{
			Path: ReadPath,
		}).Info("unmarshal request error")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//打印request相关信息
	log.Logger.WithFields(logrus.Fields{
		Path:             ReadPath,
		"len of queries": strconv.Itoa(len(request.Queries)),
	}).Info("request is " + request.String())
	//读取数据
	queryResult, err := Storage.Read(request.Queries)
	if err != nil {
		log.Logger.WithError(err).WithFields(logrus.Fields{
			Path: ReadPath,
		}).Error("read error")
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	//编码response
	response, err := prometheus.Marshal(&prompb.ReadResponse{Results: queryResult})
	if err != nil {
		log.Logger.WithError(err).WithFields(logrus.Fields{
			Path: ReadPath,
		}).Error("marshal response error")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//设置返回信息
	ctx.Header("Content-Encoding", "snappy")
	ctx.Data(http.StatusOK, "application/x-protobuf", *response)
	//打印消费时间
	consume := time.Since(begin).Seconds()
	log.Logger.WithFields(logrus.Fields{
		Path: ReadPath,
	}).Info("consume time " + strconv.FormatFloat(consume, 'f', 3, 64))
}

// Write is a controller to write metrics to storage
func Write(ctx *gin.Context) {
	begin := time.Now()
	log.Logger.WithFields(logrus.Fields{
		Path: WritePath,
	}).Info("receive request from prometheus")
	//解析request
	request := &prompb.WriteRequest{}
	if err := prometheus.Unmarshal(request, ctx.Request); err != nil {
		log.Logger.WithError(err).WithFields(logrus.Fields{
			Path: WritePath,
		}).Info("unmarshal request error")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//打印request相关信息
	log.Logger.WithFields(logrus.Fields{
		Path: WritePath,
	}).Debug("request is " + request.String())
	log.Logger.Info("len of timeSeries is " + strconv.Itoa(len(request.Timeseries)))
	//存储数据
	if err := Storage.Write(request.Timeseries); err != nil {
		log.Logger.WithError(err).WithFields(logrus.Fields{
			Path: WritePath,
		}).Error("write error")
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	//打印消费时间
	consume := time.Since(begin).Seconds()
	log.Logger.WithFields(logrus.Fields{
		Path: WritePath,
	}).Info("consume time " + strconv.FormatFloat(consume, 'f', 3, 64))
}
