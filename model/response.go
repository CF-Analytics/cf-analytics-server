package model

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseStatisticsItemData struct {
	PreDatetime    string `json:"preDatetime"`
	Datetime       string `json:"datetime"`
	FormattedTime  string `json:"formattedTime"`
	Requests       int    `json:"requests"`
	CachedRequests int    `json:"cachedRequests"`
	Bytes          int    `json:"bytes"`
	CachedBytes    int    `json:"cachedBytes"`
	Uniques        int    `json:"uniques"`
}

type ResponseStatisticsItem struct {
	ZoneName string                       `json:"zoneName"`
	Data     []ResponseStatisticsItemData `json:"data"`
}

type ResponseGraphQLHTTPRequests1hGroups struct {
	Dimensions struct {
		Datetime string `json:"datetime"`
	} `json:"dimensions"`
	Sum struct {
		Bytes          int `json:"bytes"`
		CachedBytes    int `json:"cachedBytes"`
		CachedRequests int `json:"cachedRequests"`
		Requests       int `json:"requests"`
	} `json:"sum"`
	Uniq struct {
		Uniques int `json:"uniques"`
	} `json:"uniq"`
}

type ResponseGraphQL struct {
	Data struct {
		Viewer struct {
			Zones []struct {
				HTTPRequests1hGroups []ResponseGraphQLHTTPRequests1hGroups `json:"httpRequests1hGroups"`
			} `json:"zones"`
		} `json:"viewer"`
	} `json:"data"`
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": data,
	})
}

func ResponseError(c *gin.Context, code int) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  "",
		"data": nil,
	})
}
