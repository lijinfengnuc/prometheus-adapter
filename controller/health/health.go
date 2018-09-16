// Copyright 2018 The prometheus-adapter Authors. All Rights Reserved.

// Package health defines a controller to check health
package health

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Health is a controller to return status in json format
func Health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
