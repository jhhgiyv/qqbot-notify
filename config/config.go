package config

import (
	"encoding/json"
	"log"
	"os"
)

var BotConfig JsonConfig

type JsonConfig struct {
	BotId    string          `json:"bot_id"`
	BotToken string          `json:"bot_token"`
	GuildId  string          `json:"guild_id"`
	Channels []ChannelConfig `json:"channels"`
}

type ChannelConfig struct {
	ChannelId string `json:"channel_id"`
	Type      string `json:"type"`
}

func createConfig() error {
	config := JsonConfig{BotId: "", BotToken: "", GuildId: "", Channels: []ChannelConfig{}}
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
}
