// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package flag binds and checks command-line args
package flag

import (
	"errors"
	"flag"
	"strconv"

	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/lijinfengnuc/prometheus-adapter/util/regexp"
	"github.com/sirupsen/logrus"
)

// -- Common-Type storage name
const (
	StorageES = "ElasticSearch"
)

// -- Command-line args
const (
	WebListenPort   = "web.listen-port"
	QueryMaxSize    = "query.max-size"
	AdapterFilePath = "adapter.file-path"
	AdapterName     = "adapter.name"
	MappingFilePath = "mapping.file-path"
)

// BindFlag binds and checks command-line args
func BindFlag() error {
	//绑定命令行参数
	port := *flag.String(WebListenPort, "8090", "port to listen for API")
	log.Logger.WithFields(logrus.Fields{WebListenPort: port}).Info()

	size := *flag.Int(QueryMaxSize, -1, "maximum number of records in the query")
	log.Logger.WithFields(logrus.Fields{QueryMaxSize: strconv.Itoa(size)}).Info()

	path := *flag.String(AdapterFilePath, "adapter.yaml", "path to the adapter yaml config file")
	log.Logger.WithFields(logrus.Fields{AdapterFilePath: path}).Info()

	name := *flag.String(AdapterName, StorageES, "storage service name")
	log.Logger.WithFields(logrus.Fields{AdapterName: name}).Info()

	flag.Parse()

	//校验命令行参数
	//校验web.listen-port
	if !regexp.MatchPort(port) {
		return errors.New("custom port" + port + " match false")
	}

	return nil
}
