package controller

import (
	"cf-analytics-server/model"
	"cf-analytics-server/network"
	"github.com/gin-gonic/gin"
)

func Statistics(c *gin.Context) {
	res, err := network.CfGetStatistics()
	if err != nil {
		model.ResponseError(c, 100)
		return
	}

	model.ResponseSuccess(c, res)
}
