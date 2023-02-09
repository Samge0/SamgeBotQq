package botHandler

import (
	"SamgeWxApi/cmd/qqBot/botModel"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// QqBot 当前的机器人
var QqBot *botModel.Bot

// Run 运行Bot
func Run(bot *botModel.Bot) {
	QqBot = bot
	if err := InitLogger(QqBot.Config.Loglvl); err != nil {
		panic(errors.New(fmt.Sprintf("初始化日志失败：%s", err.Error())))
	}
	defer Logger.Flush()
	for {
		func(config *botModel.Config) {
			c, _, err := ConnectWS(config)
			if err != nil {
				Logger.Error(err.Error())
				return
			}

			fmt.Println("Bot启动完成")

			Connect = c // 传出接口
			// 触发监听器 OnBotStart
			go func() {
				for _, function := range Listeners.OnBotStart {
					go function()
				}
			}()

			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					Logger.Error(err.Error())
					break // 重启bot循环 防止陷入死循环
				}
				m := make(map[string]interface{})
				if err := json.Unmarshal([]byte(message), &m); err != nil {
					Logger.Error(err.Error())
					break // 重启bot循环 防止陷入死循环
				}
				go MsgParse(m)
			}
		}(&bot.Config)
		Logger.Info("Websocket will reconnect in 5s")
		time.Sleep(5 * time.Second)
	}
}
