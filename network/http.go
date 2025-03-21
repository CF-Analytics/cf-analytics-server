package network

import (
	"bytes"
	"cf-analytics-server/config"
	"cf-analytics-server/model"
	"encoding/json"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func CfGraph(body interface{}, unmarshalVariable any) error {
	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, config.CloudflareGraphQLEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+config.CloudflareAPIToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// do nothing
		}
	}(resp.Body)

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(respBody, unmarshalVariable)

	return err
}

func CfGetStatistics() ([]model.ResponseStatisticsItem, error) {
	const timeFormatLayout = "2006-01-02T15:04:05Z"

	// UTC 时间
	before24h := time.Now().UTC().Truncate(time.Hour).Add(-24 * time.Hour)
	before24hTimeFormatted := before24h.Format(timeFormatLayout)

	query := `{
	  viewer {
		zones(filter: { zoneTag_in: $zoneTags }) {
		  httpRequests1hGroups(limit: 10000, filter: {
			datetime_gt: $datetimeStart
		  }, orderBy: [datetime_ASC]) {
			sum {
			  requests
			  cachedRequests
			  bytes
			  cachedBytes
			}
			uniq {
			  uniques
			}
			dimensions {
			  datetime
			}
		  }
		}
	  }
	}`

	reqBody := model.GraphQLRequest{
		Query: query,
		Variables: map[string]interface{}{
			"zoneTags":      config.ZoneIDs,
			"datetimeStart": before24hTimeFormatted,
		},
	}
	var result model.ResponseGraphQL

	err := CfGraph(reqBody, &result)
	if err != nil {
		return nil, err
	}

	if len(config.ZoneIDs) != len(config.ZoneNames) || len(config.ZoneIDs) != len(result.Data.Viewer.Zones) {
		return nil, errors.New("zones mismatch")
	}

	// 处理成前端可使用的数据结构
	var finalData []model.ResponseStatisticsItem

	for i, zone := range result.Data.Viewer.Zones {
		var statisticsItem model.ResponseStatisticsItem
		statisticsItem.ZoneName = config.ZoneNames[i]

		currentProcessTime := before24h.Add(1 * time.Hour)
		currentProcessTimeFormatted := currentProcessTime.Format(timeFormatLayout)

		// 循环完整的 24小时 时间间隔次数
		for timeInterval := 0; timeInterval < 24; timeInterval++ {
			var item *model.ResponseGraphQLHTTPRequests1hGroups

			// 判断 Cloudflare 是否返回了当前时间间隔的数据
			for _, group := range zone.HTTPRequests1hGroups {
				if currentProcessTimeFormatted == group.Dimensions.Datetime {
					item = &group
				}
			}

			// 将时间转为用户指定的时区时间
			currentProcessTimeLoc := currentProcessTime.In(config.GetTimeLocation())
			currentProcessPreTimeFormatted := currentProcessTimeLoc.Add(-1 * time.Hour).Format("2006/01/02 15:04")
			targetTimeFormatted := currentProcessTimeLoc.Format("2006/01/02 15:04")
			formattedTime := fmt.Sprintf("从: %s\n到: %s", currentProcessPreTimeFormatted, targetTimeFormatted)

			if item == nil {
				statisticsItem.Data = append(statisticsItem.Data, model.ResponseStatisticsItemData{
					PreDatetime:    currentProcessPreTimeFormatted,
					Datetime:       targetTimeFormatted,
					FormattedTime:  formattedTime,
					Requests:       0,
					CachedRequests: 0,
					Bytes:          0,
					CachedBytes:    0,
					Uniques:        0,
				})
			} else {
				statisticsItem.Data = append(statisticsItem.Data, model.ResponseStatisticsItemData{
					PreDatetime:    currentProcessPreTimeFormatted,
					Datetime:       targetTimeFormatted,
					FormattedTime:  formattedTime,
					Requests:       item.Sum.Requests,
					CachedRequests: item.Sum.CachedRequests,
					Bytes:          item.Sum.Bytes,
					CachedBytes:    item.Sum.CachedBytes,
					Uniques:        item.Uniq.Uniques,
				})
			}

			currentProcessTime = currentProcessTime.Add(1 * time.Hour)
			currentProcessTimeFormatted = currentProcessTime.Format(timeFormatLayout)
		}

		finalData = append(finalData, statisticsItem)
	}

	return finalData, nil
}

func PushTelegramBot() {
	bot, err := tgbotapi.NewBotAPI(config.TgBotToken)
	if err != nil {
		return
	}

	chatID, err := strconv.Atoi(config.TgUserChatID)
	if err != nil {
		return
	}

	res, err := CfGetStatistics()
	if err != nil {
		return
	}

	text := ""

	for _, item := range res {
		sumRequests := 0
		sumCachedRequests := 0
		sumBytes := 0
		sumCachedBytes := 0
		sumUniques := 0

		for _, itemData := range item.Data {
			sumRequests += itemData.Requests
			sumCachedRequests += itemData.CachedRequests
			sumBytes += itemData.Bytes
			sumCachedBytes += itemData.CachedBytes
			sumUniques += itemData.Uniques
		}

		text += fmt.Sprintf(`
		域名：%s
		近24小时统计数据：
		总请求数：%d
		已缓存的请求数：%d
		总流量：%.2f KB
		已缓存的流量：%.2f KB
		唯一访问者：%d
        `, item.ZoneName, sumRequests, sumCachedRequests, float64(sumBytes)/1000, float64(sumCachedBytes)/1000, sumUniques) + "\n\n"
	}

	_, _ = bot.Send(tgbotapi.NewMessage(int64(chatID), strings.TrimSpace(text)))
}
