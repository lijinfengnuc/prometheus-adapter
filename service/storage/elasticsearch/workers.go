// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package elasticsearch defines the storage of ES
package elasticsearch

import (
	"encoding/json"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"time"
)

type Workers []*Worker

// Worker is a struct to print
type Worker struct {
	Index        int
	Queued       int64
	LastDuration time.Duration
}

// InitWorkers converts BulkProcessorWorkerStats into Workers
func (workers *Workers) InitWorkers(bulkWorkers []*elastic.BulkProcessorWorkerStats) error {
	if len(bulkWorkers) <= 0 {
		err := "len of bulkWorkers is lte zero"
		log.Logger.Error(err)
		return errors.New(err)
	}
	for index, item := range bulkWorkers {
		*workers = append(*workers, &Worker{Index: index, Queued: item.Queued, LastDuration: item.LastDuration})
	}
	return nil
}

// String converts Workers into string
func (workers *Workers) String() string {
	jsonStr, err := json.Marshal(workers)
	if err != nil {
		log.Logger.WithError(err).Error("json marshal workers error")
		return ""
	}
	return string(jsonStr)
}
