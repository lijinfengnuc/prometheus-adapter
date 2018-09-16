// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package flag defines some utils func about flag
package flag

import (
	"flag"
	"strconv"
	"time"
)

// GetFlag gets value in string format from flag
func GetFlag(name string) *string {
	flagPoint := flag.Lookup(name)
	if flagPoint == nil {
		return nil
	}
	if flagPoint.Value == nil {
		return nil
	}
	value := flagPoint.Value.String()
	return &value
}

// GetStringFlag gets value in string format from flag
func GetStringFlag(name string) *string {
	valueStr := GetFlag(name)
	if valueStr != nil {
		return valueStr
	} else {
		return nil
	}
}

// GetBoolFlag gets value in bool format from flag
func GetBoolFlag(name string) *bool {
	valueStr := GetFlag(name)
	if valueStr != nil {
		value, err := strconv.ParseBool(*valueStr)
		if err != nil {
			return nil
		}
		return &value
	} else {
		return nil
	}
}

// GetIntFlag gets value in int format from flag
func GetIntFlag(name string) *int {
	valueStr := GetFlag(name)
	if valueStr != nil {
		value, err := strconv.Atoi(*valueStr)
		if err != nil {
			return nil
		}
		return &value
	} else {
		return nil
	}
}

// GetInt64Flag gets value in int64 format from flag
func GetInt64Flag(name string) *int64 {
	valueStr := GetFlag(name)
	if valueStr != nil {
		value, err := strconv.ParseInt(*valueStr, 10, 64)
		if err != nil {
			return nil
		}
		return &value
	} else {
		return nil
	}
}

// GetFloat64Flag gets value in float64 format from flag
func GetFloat64Flag(name string) *float64 {
	valueStr := GetFlag(name)
	if valueStr != nil {
		value, err := strconv.ParseFloat(*valueStr, 64)
		if err != nil {
			return nil
		}
		return &value
	} else {
		return nil
	}
}

// GetDurationFlag gets value in duration format from flag
func GetDurationFlag(name string) *time.Duration {
	valueStr := GetFlag(name)
	if valueStr != nil {
		value, err := time.ParseDuration(*valueStr)
		if err != nil {
			return nil
		}
		return &value
	} else {
		return nil
	}
}
