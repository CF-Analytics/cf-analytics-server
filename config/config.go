package config

import (
	"context"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"log"
	"os"
	"strings"
	"time"
)

const CloudflareGraphQLEndpoint = "https://api.cloudflare.com/client/v4/graphql"

var (
	Debug              = os.Getenv("CFA_DEBUG") == "true"
	CloudflareAPIToken = os.Getenv("CLOUDFLARE_API_TOKEN")
	TgBotToken         = os.Getenv("TG_BOT_TOKEN")
	TgUserChatID       = os.Getenv("TG_USER_CHAT_ID")
	ZoneIDs            []string
	ZoneNames          []string

	// TimeLocation Asia/Shanghai
	timeLocation *time.Location
)

func Init() {
	userZoneIDs := os.Getenv("CLOUDFLARE_ZONE_ID")

	// 判断环境变量是否存在并且有值
	if len(CloudflareAPIToken) <= 0 || len(userZoneIDs) <= 0 {
		log.Fatalln("缺少必要的环境变量")
	}

	// 加载用户设置的时区
	userTimeLocation := os.Getenv("CFA_TIME_LOCATION")
	if len(userTimeLocation) <= 0 {
		userTimeLocation = time.Local.String()
	}

	var err error
	timeLocation, err = time.LoadLocation(userTimeLocation)
	if err != nil {
		log.Fatalln("时区加载失败")
	}

	// 处理 域 数据
	client := cloudflare.NewClient(option.WithAPIToken(CloudflareAPIToken))
	page, err := client.Zones.List(context.TODO(), zones.ZoneListParams{PerPage: cloudflare.Float(50)})
	if err != nil {
		log.Fatalln("Cloudflare 域数据获取失败")
	}

	ZoneIDs = strings.Split(userZoneIDs, ",")
	for _, zoneID := range ZoneIDs {
		for _, zone := range page.Result {
			if zoneID == zone.ID {
				ZoneNames = append(ZoneNames, zone.Name)
			}
		}
	}
}

func GetTimeLocation() *time.Location {
	return timeLocation
}
