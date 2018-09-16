// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package elasticsearch defines the storage of ES
package elasticsearch

import (
	"context"
	"io"
	"reflect"
	"strconv"

	"encoding/json"
	"github.com/lijinfengnuc/prometheus-adapter/flag"
	flagUtil "github.com/lijinfengnuc/prometheus-adapter/util/flag"
	jsonUtil "github.com/lijinfengnuc/prometheus-adapter/util/json"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/lijinfengnuc/prometheus-adapter/util/os/path"
	"github.com/lijinfengnuc/prometheus-adapter/util/regexp"
	"github.com/lijinfengnuc/prometheus-adapter/util/yaml"
	"github.com/olivere/elastic"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/prompb"
	"github.com/sirupsen/logrus"
)

// -- Some constants
const (
	Index      = "index"
	Type       = "type"
	QueryIndex = "queryIndex"
	Page       = "page"
)

// ElasticCluster defines some fields about ES cluster
type ElasticCluster struct {
	ElasticNodes []*ElasticNode `yaml:"elasticNodes"`
	User         string         `yaml:"user"`
	Password     string         `yaml:"password"`
	Index        string         `yaml:"index"`
	TypeAlias    string         `yaml:"type"`
	Sniff        bool           `yaml:"sniff"`
	Healthcheck  bool           `yaml:"healthcheck"`
	Workers      int            `yaml:"workers"`
	BulkSize     int            `yaml:"bulkSize"`
	QuerySize    int            `yaml:"querySize"`
	MappingPath  string         `yaml:"mappingPath"`
	Client       *elastic.Client
}

// ElasticNode defines some fields about ES node
type ElasticNode struct {
	IP   string `yaml:"ip"`
	Port string `yaml:"port"`
}

// loadConfig loads fields from file into struct ElasticCluster and checks them
func (elasticCluster *ElasticCluster) loadConfig() error {
	//获取配置文件路径
	adapterFilePath, err := path.GetPath(*flagUtil.GetStringFlag(flag.AdapterFilePath))
	if err != nil {
		log.Logger.Error("get adapter file path error")
		return err
	}
	log.Logger.WithFields(logrus.Fields{
		flag.AdapterFilePath: adapterFilePath,
	}).Info("get adapter file path success")

	//加载对应storage的配置文件
	if err := yaml.Unmarshal(elasticCluster, adapterFilePath); err != nil {
		log.Logger.WithFields(logrus.Fields{
			flag.AdapterFilePath: adapterFilePath,
		}).Error("load adapter file error")
		return err
	}
	log.Logger.WithFields(logrus.Fields{
		flag.AdapterFilePath: adapterFilePath,
	}).Info("load adapter file success")

	//校验字段
	//校验ElasticNode IP/Port
	if len(elasticCluster.ElasticNodes) == 0 {
		elasticCluster.ElasticNodes = []*ElasticNode{{IP: "127.0.0.1", Port: "9200"}}
	}
	for _, elasticNode := range elasticCluster.ElasticNodes {
		//检查ip,port格式
		if !regexp.MatchIp(elasticNode.IP) {
			return errors.New("ip " + elasticNode.IP + " is not match pattern")
		} else if !regexp.MatchPort(elasticNode.Port) {
			return errors.New("port" + elasticNode.Port + " is not match pattern")
		}
		log.Logger.Info(elasticNode.IP + ":" + elasticNode.Port + " join cluster")
	}
	//校验user
	if elasticCluster.User == "" {
		elasticCluster.User = "elastic"
	}
	log.Logger.WithFields(logrus.Fields{"user": elasticCluster.User}).Info()
	//校验password
	if elasticCluster.Password == "" {
		elasticCluster.Password = "changeme"
	}
	log.Logger.WithFields(logrus.Fields{"password": elasticCluster.Password}).Info()
	//校验index
	if elasticCluster.Index == "" {
		elasticCluster.Index = "prometheus"
	}
	log.Logger.WithFields(logrus.Fields{"index": elasticCluster.Index}).Info()
	//校验type
	if elasticCluster.TypeAlias == "" {
		elasticCluster.TypeAlias = "metric"
	}
	log.Logger.WithFields(logrus.Fields{"type": elasticCluster.TypeAlias}).Info()
	//校验sniff,初始化默认false
	log.Logger.WithFields(logrus.Fields{"sniff": elasticCluster.Sniff}).Info()
	//校验healthcheck,初始化默认false
	log.Logger.WithFields(logrus.Fields{"healthcheck": elasticCluster.Healthcheck}).Info()
	//校验workers
	if elasticCluster.Workers == 0 {
		elasticCluster.Workers = 1
	}
	log.Logger.WithFields(logrus.Fields{"workers": strconv.Itoa(elasticCluster.Workers)}).Info()
	//校验bulkSize
	if elasticCluster.BulkSize == 0 {
		elasticCluster.BulkSize = 1
	}
	log.Logger.WithFields(logrus.Fields{"bulkSize": strconv.Itoa(elasticCluster.BulkSize)}).Info()
	//校验querySize
	if elasticCluster.QuerySize == 0 {
		elasticCluster.QuerySize = 5000
	} else if elasticCluster.QuerySize > 10000 {
		return errors.New(adapterFilePath + ":querySize should less than 10000")
	}
	log.Logger.WithFields(logrus.Fields{"querySize": strconv.Itoa(elasticCluster.QuerySize)}).Info()
	//校验mappingPath
	if elasticCluster.MappingPath == "" {
		elasticCluster.MappingPath = "mapping.json"
	}
	log.Logger.WithFields(logrus.Fields{"mappingPath": elasticCluster.MappingPath}).Info()

	return nil
}

// Init implements Init method of interface Storage
func (elasticCluster *ElasticCluster) Init() error {
	//加载配置文件
	if err := elasticCluster.loadConfig(); err != nil {
		log.Logger.Error("load adapter file error")
		return err
	}
	log.Logger.Info("load adapter file success")

	//拼接urls
	var urls []string
	for _, elasticNode := range elasticCluster.ElasticNodes {
		url := "http://" + elasticNode.IP + ":" + elasticNode.Port
		urls = append(urls, url)
	}

	//创建client
	elasticClient, err := elastic.NewClient(elastic.SetURL(urls...), elastic.SetSniff(elasticCluster.Sniff),
		elastic.SetHealthcheck(elasticCluster.Healthcheck), elastic.SetBasicAuth(elasticCluster.User, elasticCluster.Password))
	if err != nil {
		log.Logger.Error("create client for ES error")
		return err
	}
	log.Logger.Info("create client for ES success")

	//client赋值
	elasticCluster.Client = elasticClient

	//验证index/type是否存在
	if mapping, err := elasticClient.GetMapping().Index(elasticCluster.Index).
		Type(elasticCluster.TypeAlias).Do(context.Background()); err != nil || len(mapping) == 0 {
		log.Logger.Warn("type is not exist")
		if err := elasticCluster.createType(); err != nil {
			log.Logger.Error("create type error")
			return err
		}
		log.Logger.Info("create type success")
	} else {
		log.Logger.Info("type is already exist")
	}
	return nil
}

// createType creates specific index\type in ES
func (elasticCluster *ElasticCluster) createType() error {

	mappingPath, err := path.GetPath(elasticCluster.MappingPath)
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			Index: elasticCluster.Index,
		}).Error("get mapping path error")
		return err
	}
	//加载mapping file
	var mapping map[string]interface{}
	if err := jsonUtil.Unmarshal(&mapping, mappingPath); err != nil {
		log.Logger.WithFields(logrus.Fields{
			Index: elasticCluster.Index,
		}).Error("unmarshal mapping file error")
		return err
	}

	//client赋值
	client := elasticCluster.Client

	//检测index是否存在
	indexExist, err := client.IndexExists().Index([]string{elasticCluster.Index}).Do(context.Background())
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			Index: elasticCluster.Index,
		}).Error("check index exist error")
		return err
	}

	if !indexExist {
		//index不存在创建index
		log.Logger.WithFields(logrus.Fields{
			Index: elasticCluster.Index,
		}).Warn("index not exist,create...")
		result, err := client.CreateIndex(elasticCluster.Index).Do(context.Background())
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				Index: elasticCluster.Index,
			}).Error("create index error")
			return err
		}
		if result.Acknowledged && result.ShardsAcknowledged {
			log.Logger.WithFields(logrus.Fields{
				Index: elasticCluster.Index,
			}).Info("create index success")
		} else {
			log.Logger.WithFields(logrus.Fields{
				Index: elasticCluster.Index,
			}).Error("create index error")
			return errors.New("Acknowledged or ShardsAcknowledged is false when create index")
		}
	} else {
		//index存在，打印日志
		log.Logger.WithFields(logrus.Fields{
			Index: elasticCluster.Index,
		}).Info("index already exist")
	}
	//创建type
	result, err := client.PutMapping().Index(elasticCluster.Index).Type(elasticCluster.TypeAlias).
		BodyJson(mapping).Do(context.Background())
	if err != nil {
		log.Logger.WithFields(logrus.Fields{
			Index: elasticCluster.Index,
		}).Error("put mapping error")
		return err
	}
	if result.Acknowledged {
		log.Logger.WithFields(logrus.Fields{
			Index: elasticCluster.Index,
		}).Info("put mapping success")
	} else {
		log.Logger.WithFields(logrus.Fields{
			Index: elasticCluster.Index,
		}).Error("put mapping error")
		return errors.New("Acknowledged is false when put mapping")
	}

	return nil
}

// Write implements Write method of interface Storage
func (elasticCluster *ElasticCluster) Write(timeSeries []*prompb.TimeSeries) error {
	client := elasticCluster.Client
	//创建BulkProcessor
	bulkProcessor, err := client.BulkProcessor().Workers(elasticCluster.Workers).
		BulkSize(elasticCluster.BulkSize << 20).Stats(true).After(after).Do(context.Background())
	if err != nil {
		log.Logger.Error("create BulkProcessor error")
		return err
	}
	log.Logger.Info("create BulkProcessor success")

	//关闭BulkProcessor
	defer bulkProcessor.Close()

	//循环构建sample并存储
	var samples Samples
	samples.TimeSeries2Samples(timeSeries)
	//循环存储
	for _, sample := range samples {
		//创建BulkIndexRequest
		bulkRequest := elastic.NewBulkIndexRequest().Index(elasticCluster.Index).
			Type(elasticCluster.TypeAlias).Doc(sample)
		//存储
		bulkProcessor.Add(bulkRequest)
	}

	//清空管道
	if err := bulkProcessor.Flush(); err != nil {
		log.Logger.Error("flush for last commit error")
		stats(bulkProcessor.Stats())
		return err
	}
	log.Logger.Info("flush for last commit success")

	//打印执行日志信息
	stats(bulkProcessor.Stats())

	return nil
}

// after prints commit detail after every commit
func after(executionId int64, requests []elastic.BulkableRequest, response *elastic.BulkResponse, err error) {
	if err != nil {
		//打印错误日志
		log.Logger.WithError(err).Error("executionId:" +
			strconv.FormatInt(executionId, 10) + " commit error")

		//循环打印每个bulkRequest的source
		for index, bulkRequest := range requests {
			source, err := bulkRequest.Source()
			if err != nil {
				log.Logger.WithError(err).WithFields(logrus.Fields{
					"index": index,
				}).Error("bulkRequest source error")
			} else {
				log.Logger.WithFields(logrus.Fields{
					"index":  index,
					"source": source,
				}).Debug("bulkRequest source success")
			}
		}

		//循环打印每个response
		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Logger.WithError(err).Error("marshal response error")
		} else {
			log.Logger.WithError(err).WithFields(logrus.Fields{
				"response": string(responseBytes),
			}).Error("response detail")
		}

		//后期可于此处添加错误处理机制
		//1.停止BulkProcessor
		//2.存储source
	}
}

// stats prints execute detail after all commit
func stats(stats elastic.BulkProcessorStats) {
	//构建Workers的json文本
	var workersStr string
	var workers Workers
	if err := workers.InitWorkers(stats.Workers); err == nil {
		workersStr = workers.String()
	}

	//打印stats信息
	log.Logger.WithFields(logrus.Fields{
		"Flushed":   stats.Flushed,
		"Committed": stats.Committed,
		"Indexed":   stats.Indexed,
		"Created":   stats.Created,
		"Updated":   stats.Updated,
		"Deleted":   stats.Deleted,
		"Succeeded": stats.Succeeded,
		"Failed":    stats.Failed,
		"Workers":   workersStr,
	}).Info("stats info detail")
}

// Read implements Read method of interface Storage
func (elasticCluster *ElasticCluster) Read(queries []*prompb.Query) ([]*prompb.QueryResult, error) {
	//初始化queryResult
	queryResults := make([]*prompb.QueryResult, 0, len(queries))

	//循环查询（应该只有一个元素，不循环也行）
	for index, query := range queries {
		log.Logger.WithFields(logrus.Fields{
			QueryIndex: index,
		}).Info("query start")

		//新建组合查询条件
		boolQuery, err := buildBoolQuery(query)
		if err != nil {
			log.Logger.WithError(err).WithFields(logrus.Fields{
				QueryIndex: index,
			}).Error("build BoolQuery error")
			continue
		}

		//根据查询条件分页查询查询
		samples, err := elasticCluster.scrollSaerch(boolQuery)
		if err != nil {
			log.Logger.WithError(err).WithFields(logrus.Fields{
				QueryIndex: index,
			}).Error("scroll search error")
			continue
		}
		log.Logger.Info("scroll search success")

		//将查询结果转化为samples
		if samples != nil {
			queryResult := samples.Samples2QueryResult()
			queryResults = append(queryResults, queryResult)
		} else {
			log.Logger.WithFields(logrus.Fields{
				QueryIndex: index,
			}).Info("count is 0")
		}

		log.Logger.WithFields(logrus.Fields{
			QueryIndex: index,
		}).Info("query end")
	}

	return queryResults, nil
}

// buildBoolQuery builds a bool query for query
func buildBoolQuery(query *prompb.Query) (*elastic.BoolQuery, error) {
	elastic.NewNestedAggregation()
	boolQuery := elastic.NewBoolQuery()
	//标签过滤
	for _, matcher := range query.Matchers {
		switch matcher.Type {
		case prompb.LabelMatcher_EQ:
			boolQuery.Must(elastic.NewTermQuery("labels."+matcher.Name+".keyword", matcher.Value))
		case prompb.LabelMatcher_NEQ:
			boolQuery.MustNot(elastic.NewTermQuery("labels."+matcher.Name+".keyword", matcher.Value))
		case prompb.LabelMatcher_RE:
			matcher.Value = regexp.RevisePattern(matcher.Value)
			boolQuery.Must(elastic.NewRegexpQuery("labels."+matcher.Name+".keyword", matcher.Value))
		case prompb.LabelMatcher_NRE:
			matcher.Value = regexp.RevisePattern(matcher.Value)
			boolQuery.MustNot(elastic.NewRegexpQuery("labels."+matcher.Name+".keyword", matcher.Value))
		default:
			return nil, errors.New("matcher type " + matcher.Type.String() + " not match any case")
		}
	}
	//时间过滤
	boolQuery.Filter(elastic.NewRangeQuery("timestamp").Gte(query.StartTimestampMs).Lte(query.EndTimestampMs))
	return boolQuery, nil
}

// scrollSaerch queries metrics by page
func (elasticCluster *ElasticCluster) scrollSaerch(boolQuery *elastic.BoolQuery) (*Samples, error) {
	var samples Samples
	var count int

	//查询总数
	scrollService := elasticCluster.Client.Scroll().KeepAlive("3m").Index(elasticCluster.Index).
		Type(elasticCluster.TypeAlias).Query(boolQuery).Size(elasticCluster.QuerySize).
		Sort("timestamp", true)

	//关闭service
	defer scrollService.Clear(context.Background())

	//分页查询
	for page := 1; true; page++ {
		//查询
		pageResult, err := scrollService.Do(context.Background())
		//count为0
		if err == io.EOF && page == 1 {
			return nil, nil
		}
		//错误处理
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				Page: page,
			}).Error("page search error")
			return nil, err
		}
		//设置count
		if page == 1 {
			count = int(pageResult.Hits.TotalHits)
			queryMaxSize := *flagUtil.GetIntFlag(flag.QueryMaxSize)
			if queryMaxSize > 0 && count > queryMaxSize {
				count = queryMaxSize
			}
			log.Logger.Info("count is " + strconv.Itoa(count))
		}
		//遍历获取查询结果
		for _, sample := range pageResult.Each(reflect.TypeOf(Sample{})) {
			sample := sample.(Sample)
			samples = append(samples, &sample)
		}
		log.Logger.WithFields(logrus.Fields{
			Page: page,
		}).Info("page search success")
		//结束循环
		if page*elasticCluster.QuerySize >= count {
			break
		}
	}
	//若结果集大于最大长度，则截取
	if len(samples) > count {
		samples = samples[0:count]
	}

	return &samples, nil
}
