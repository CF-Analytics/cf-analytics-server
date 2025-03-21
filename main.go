package main

import (
	"cf-analytics-server/config"
	"cf-analytics-server/network"
	"cf-analytics-server/route"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	config.Init()

	if len(config.TgBotToken) > 0 && len(config.TgUserChatID) > 0 {
		c := cron.New(cron.WithLocation(config.GetTimeLocation()))

		// 每天 00:10 执行
		_, err := c.AddFunc("10 0 * * *", network.PushTelegramBot)
		if err != nil {
			log.Fatalln("Telegram 定时推送任务添加失败")
		}

		c.Start()
		defer c.Stop()
	}

	route.Start()
}
