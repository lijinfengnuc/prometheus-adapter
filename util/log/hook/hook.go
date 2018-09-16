// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package hook defines a common struct for inheritance
package hook

import "github.com/sirupsen/logrus"

// Hook is a common struct for inheritance
type Hook struct {
	Level []logrus.Level
}
