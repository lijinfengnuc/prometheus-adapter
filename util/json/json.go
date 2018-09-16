// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package json defines some utils about json
package json

import (
	"encoding/json"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const (
	JsonFilePath = "json-file-path"
)

// Unmarshal converts json file into struct
func Unmarshal(jsonStruct interface{}, filePath string) error {
	//读取json文件
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			JsonFilePath: filePath,
		}).Error("read file error")
		return err
	}

	//转化成对应的结构体
	if err := json.Unmarshal(file, jsonStruct); err != nil {
		log.Logger.WithFields(logrus.Fields{
			JsonFilePath: filePath,
		}).Error("unmarshal file error")
		return err
	}
	return nil
}
