// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package storage defines a interface Storage and a func GetStorage
package storage

import (
	"github.com/lijinfengnuc/prometheus-adapter/flag"
	"github.com/lijinfengnuc/prometheus-adapter/service/storage/elasticsearch"
	flagUtil "github.com/lijinfengnuc/prometheus-adapter/util/flag"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/prompb"
)

// Storage defines some method as a common storage
type Storage interface {
	Init() error
	Write(timeSeries []*prompb.TimeSeries) error
	Read(queries []*prompb.Query) ([]*prompb.QueryResult, error)
}

// GetStorage returns a specific storage
func GetStorage() (Storage, error) {
	var storage Storage
	//动态创建storage
	adapterName := *flagUtil.GetStringFlag(flag.AdapterName)
	switch adapterName {
	case flag.StorageES:
		storage = &elasticsearch.ElasticCluster{}
	default:
		return nil, errors.New("storage name " + adapterName + " not match any case")
	}

	//初始化storage
	if err := storage.Init(); err != nil {
		return nil, err
	}

	return storage, nil
}
