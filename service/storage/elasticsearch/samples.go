// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package elasticsearch defines the storage of ES
package elasticsearch

import (
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/prompb"
	"math"
)

type Samples []*Sample

// Sample is struct for saving in ES
type Sample struct {
	Labels    model.Metric `json:"labels"`
	Value     float64      `json:"value"`
	TimeStamp int64        `json:"timestamp"`
}

// TimeSeries2Samples converts TimeSeries into Samples
func (samples *Samples) TimeSeries2Samples(timeSeries []*prompb.TimeSeries) {
	for _, ts := range timeSeries {
		//构建Metric
		metric := make(model.Metric, len(ts.Labels))
		for _, label := range ts.Labels {
			metric[model.LabelName(label.Name)] = model.LabelValue(label.Value)
		}

		//构建samples
		for _, sample := range ts.Samples {
			if math.IsNaN(sample.Value) {
				sample.Value = 0
			}
			*samples = append(*samples,
				&Sample{metric, sample.Value, sample.Timestamp})
		}
	}
}

// Samples2QueryResult converts Samples into TimeSeries
func (samples *Samples) Samples2QueryResult() *prompb.QueryResult {
	timeSeriesMap := make(map[string]*prompb.TimeSeries)
	for _, sample := range *samples {
		//获取指标指纹
		fingerprint := sample.Labels.Fingerprint().String()

		//获取指标的ts
		ts, ok := timeSeriesMap[fingerprint]
		if !ok {
			//timeSeries中没有对应指标，初始化ts
			labels := make([]*prompb.Label, 0, len(sample.Labels))
			//构建labels
			for name, value := range sample.Labels {
				labels = append(labels,
					&prompb.Label{Name: string(name), Value: string(value)})
			}
			//初始化samples
			sps := make([]*prompb.Sample, 0, len(*samples))
			//TimeSeries初始化
			timeSeriesMap[fingerprint] = &prompb.TimeSeries{Labels: labels, Samples: sps}
			ts = timeSeriesMap[fingerprint]
		}

		//构建samples
		ts.Samples = append(ts.Samples,
			&prompb.Sample{Value: sample.Value, Timestamp: sample.TimeStamp})
	}
	timeSeries := make([]*prompb.TimeSeries, 0, len(timeSeriesMap))
	for _, ts := range timeSeriesMap {
		timeSeries = append(timeSeries, ts)
	}

	//后期有需要根据ReadHints.StepMs进行筛选，暂时没想好怎么做
	return &prompb.QueryResult{Timeseries: timeSeries}
}
