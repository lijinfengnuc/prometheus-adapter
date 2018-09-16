// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package storage defines Read/Write controller for prometheus
package storage

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/lijinfengnuc/prometheus-adapter/util/json"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/lijinfengnuc/prometheus-adapter/util/prometheus"
	"github.com/prometheus/prometheus/prompb"
	"io/ioutil"
	"net/http"
	"testing"
)

// -- Some constants
const (
	WriteFilePath = "write_test.json"
	ReadFilePath  = "read_test.json"
	Url           = "http://127.0.0.1:8090"
)

// TestWrite tests Write controller
func TestWrite(t *testing.T) {
	var timeSeries []*prompb.TimeSeries
	if err := json.Unmarshal(&timeSeries, WriteFilePath); err != nil {
		log.Logger.WithError(err).Error("unmarshal " + WriteFilePath + " err")
		return
	}
	writeData, err := prometheus.Marshal(&prompb.WriteRequest{Timeseries: timeSeries})
	if err != nil {
		log.Logger.WithError(err).Error("marshal write data err")
		return
	}
	client := &http.Client{}
	resp, err := client.Post(Url+"/v1/write", "application/x-protobuf", bytes.NewReader(*writeData))
	if err != nil {
		log.Logger.WithError(err).Error("request write controller err")
		return
	}
	log.Logger.Info(resp)
}

// TestRead tests Read controller
func TestRead(t *testing.T) {
	var queries []*prompb.Query
	if err := json.Unmarshal(&queries, ReadFilePath); err != nil {
		log.Logger.WithError(err).Error("unmarshal " + ReadFilePath + " err")
		return
	}
	readQuery, err := prometheus.Marshal(&prompb.ReadRequest{Queries: queries})
	if err != nil {
		log.Logger.WithError(err).Error("marshal read query err")
		return
	}
	client := &http.Client{}
	response, err := client.Post(Url+"/v1/read", "application/x-protobuf", bytes.NewReader(*readQuery))
	if err != nil {
		log.Logger.WithError(err).Error("request read controller err")
		return
	}
	readResponse := &prompb.ReadResponse{}
	if err := Unmarshal(readResponse, response); err != nil {
		log.Logger.WithError(err).Error("unmarshal response err")
		return
	}
	log.Logger.Info(readResponse.String())
}

func Unmarshal(message proto.Message, response *http.Response) error {
	//读取request
	compressedBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Logger.Error("read compressedBody error")
		return err
	}
	//解压request
	body, err := snappy.Decode(nil, compressedBody)
	if err != nil {
		log.Logger.Error("decode body error")
		return err
	}
	//解码request
	if err := proto.Unmarshal(body, message); err != nil {
		log.Logger.Error("unmarshal request error")
		return err
	}
	return nil
}
