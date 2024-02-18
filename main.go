package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jhhgiyv/qqbot-notify/config"
	_ "github.com/jhhgiyv/qqbot-notify/config"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/errs"
	"github.com/tencent-connect/botgo/log"
	"regexp"
	"strings"
)

type payload struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func main() {
	route := gin.Default()
	route.POST("/", notify)
	err := route.Run(config.BotConfig.Addr)
	if err != nil {
		log.Error(err)
	}
}

func filterURL(message string) string {
	urlRegexp := regexp.MustCompile(`https?://[^\s]+`)
	if urlRegexp.MatchString(message) {
		url := urlRegexp.FindString(message)
		url = strings.ReplaceAll(url, "http://", "")
		url = strings.ReplaceAll(url, "https://", "")
		url = strings.ReplaceAll(url, "/", "\\")
		url = strings.ReplaceAll(url, ".", "。")
		message = urlRegexp.ReplaceAllString(message, url)
	}
	return message
}

func notify(context *gin.Context) {
	payloadObj := payload{}
	err := context.ShouldBindJSON(&payloadObj)
	if err != nil {
		context.JSON(400, gin.H{"code": 400, "msg": err.Error()})
	}
	s := strings.Split(payloadObj.Subject, "/")
	var content string
	var ok bool
	var channelId string
	channelId = config.BotConfig.Channels[0].ChannelId
	if len(s) > 1 {
		for _, channel := range config.BotConfig.Channels {
			group := channel.Group
			if s[1] == group {
				channelId = channel.ChannelId
				content = fmt.Sprintf("%s\n%s", s[0], payloadObj.Message)
				ok = true
				break
			}
		}
	}
	if !ok {
		message := filterURL(payloadObj.Message)
		content = fmt.Sprintf("%s\n%s", payloadObj.Subject, message)
	}
	log.Infof("发送:\n\"\"\"\n%s\n\"\"\"", content)
	_, err = config.Api.PostMessage(config.Ctx, channelId, &dto.MessageToCreate{Content: content})
	if err != nil {
		var er *errs.Err
		ok = errors.As(err, &er)
		if ok && er.Code() == 202 {
			return
		}
		log.Error(err)
		return
	}
}
