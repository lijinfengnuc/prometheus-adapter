// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Prometheus package defines some utils about prometheus
package prometheus

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"io/ioutil"
	"net/http"
)

// Unmarshal converts http-request into proto-Message
func Unmarshal(message proto.Message, request *http.Request) error {
	//读取request
	compressedBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Logger.Error("read request error")
		return err
	}

	//解压request
	body, err := snappy.Decode(nil, compressedBody)
	if err != nil {
		log.Logger.Error("decode compressedBody error")
		return err
	}

	//解码request
	if err := proto.Unmarshal(body, message); err != nil {
		log.Logger.Error("unmarshal body error")
		return err
	}

	return nil
}

// Marshal converts proto-Message into *[]byte
func Marshal(message proto.Message) (*[]byte, error) {
	//编码response
	response, err := proto.Marshal(message)
	if err != nil {
		log.Logger.Error("marshal proto message error")
		return nil, err
	}

	//压缩response
	compressedResponse := snappy.Encode(nil, response)

	return &compressedResponse, nil
}
