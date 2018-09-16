// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package regexp defines some utils about regexp
package regexp

import (
	"github.com/lijinfengnuc/prometheus-adapter/util/log"
	"github.com/sirupsen/logrus"
	"regexp"
)

const (
	Pattern = "pattern"
)

// RevisePattern removes "^" and "$" in pattern
func RevisePattern(pattern string) string {
	if pattern[0:1] == "^" {
		pattern = pattern[1:]
	}
	if pattern[len(pattern)-1:] == "$" {
		pattern = pattern[0 : len(pattern)-1]
	}
	return pattern
}

// MatchIp uses regexp to match IP
func MatchIp(ip string) bool {
	pattern := "^((2[0-4]\\d|25[0-5]|[01]?\\d\\d?)\\.){3}(2[0-4]\\d|25[0-5]|[01]?\\d\\d?)$"
	isMatch, err := regexp.MatchString(pattern, ip)
	if isMatch == false {
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				Pattern: pattern,
			}).WithError(err).Error("ip pattern is not correct")
			return false
		}
		log.Logger.WithFields(logrus.Fields{
			"ip":    ip,
			Pattern: pattern,
		}).Error("ip match false")
		return false
	}
	return true
}

// MatchPort uses regexp to match port
func MatchPort(port string) bool {
	pattern := "^([0-9]|[1-9]\\d{1,3}|[1-5]\\d{4}|6[0-5]{2}[0-3][0-5])$"
	isMatch, err := regexp.MatchString(pattern, port)
	if isMatch == false {
		if err != nil {
			log.Logger.WithFields(logrus.Fields{
				Pattern: pattern,
			}).WithError(err).Error("port pattern is not correct")
			return false
		}
		log.Logger.WithFields(logrus.Fields{
			"port":  port,
			Pattern: pattern,
		}).Error("port match false")
		return false
	}
	return true
}
