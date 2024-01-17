package config

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"log"
	"os"
	"strconv"
	"time"
)

var BotConfig JsonConfig

type JsonConfig struct {
	BotId       string          `json:"bot_id"`
	BotApiToken string          `json:"bot_api_token"`
	GuildId     string          `json:"guild_id"`
	IsSandbox   bool            `json:"is_sandbox"`
	Addr        string          `json:"addr"`
	Channels    []ChannelConfig `json:"channels"`
}

type ChannelConfig struct {
	ChannelId string `json:"channel_id"`
	Group     string `json:"group"`
}

var Api openapi.OpenAPI
var Ctx context.Context

func createConfig() error {
	config := JsonConfig{BotId: "", BotApiToken: "", GuildId: "", Channels: []ChannelConfig{}, IsSandbox: true}
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}
	file, err := os.Create("config.json")
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil

}

func init() {
	for {
		file, err := os.ReadFile("config.json")
		if err != nil {
			if os.IsNotExist(err) {
				err := createConfig()
				if err != nil {
					log.Fatal("createConfig error: ", err.Error())
				}
				continue
			} else {
				log.Fatal("InitConfig error: ", err.Error())
			}
		}
		config := JsonConfig{}
		err = json.Unmarshal(file, &config)
		if err != nil {
			log.Fatal("load config error: ", err.Error())
		}
		BotConfig = config
		break
	}
	checkConfig()
}

func printGuilds(api openapi.OpenAPI, ctx context.Context) {
	if BotConfig.GuildId != "" {
		return
	}
	guilds, err := api.MeGuilds(ctx, &dto.GuildPager{})
	if err != nil {
		log.Fatalln("printGuilds error: ", err)
	}
	fmt.Println()
	fmt.Println("未填写 guild_id(频道id, 非客户端显示的id)")
	fmt.Println("将打印机器人加入的频道")
	for _, guild := range guilds {
		fmt.Println()
		fmt.Println("频道id:", guild.ID)
		fmt.Println("频道名称:", guild.Name)
		fmt.Println("频道描述:", guild.Desc)
		fmt.Println("频道图标:", guild.Icon)
		fmt.Println("频道成员数:", guild.MemberCount)
		fmt.Println()
	}
	log.Fatalln("程序退出，请填写后再次运行")
}

func printChannels(api openapi.OpenAPI, ctx context.Context) {
	if len(BotConfig.Channels) != 0 {
		return
	}
	channels, err := api.Channels(ctx, BotConfig.GuildId)
	if err != nil {
		log.Fatalln("printChannels error: ", err)
	}
	fmt.Println("未填写 channels")
	fmt.Println("将打印频道的所有子频道")
	for _, channel := range channels {
		fmt.Println()
		fmt.Println("子频道id:", channel.ID)
		fmt.Println("子频道名称:", channel.Name)
		fmt.Println("子频道类型:", channel.Type)
		fmt.Println()
	}
	log.Fatalln("程序退出，请填写后再次运行")
}

func checkConfig() {
	if BotConfig.BotId == "" || BotConfig.BotApiToken == "" {
		log.Fatalln("未填写bot_id或bot_api_token\n请根据 https://q.qq.com/qqbot/#/developer/developer-setting 上的内容填写")
	}
	botId, err := strconv.ParseUint(BotConfig.BotId, 10, 64)
	if err != nil {
		log.Fatalln("botId error", err)
	}
	botToken := token.BotToken(botId, BotConfig.BotApiToken)
	if BotConfig.IsSandbox {
		Api = botgo.NewSandboxOpenAPI(botToken).WithTimeout(3 * time.Second)
	} else {
		Api = botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	}
	Ctx = context.Background()
	printGuilds(Api, Ctx)
	printChannels(Api, Ctx)
}
