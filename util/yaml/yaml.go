// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package yaml defines some utils about yaml
package yaml

import (
	"github.com/go-yaml/yaml"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const (
	YamlFilePath = "yaml-file-path"
)

// Unmarshal converts yaml file into struct
func Unmarshal(yamlStruct interface{}, filePath string) error {
	//读取yaml文件
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			YamlFilePath: filePath,
		}).Error("read file error")
		return err
	}

	//转化成对应的结构体
	if err := yaml.Unmarshal(file, yamlStruct); err != nil {
		log.Logger.WithFields(logrus.Fields{
			YamlFilePath: filePath,
		}).Error("unmarshal file error")
		return err
	}

	return nil
}
