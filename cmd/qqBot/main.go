package main

import (
	"SamgeWxApi/cmd/qqBot/botHandler"
	"SamgeWxApi/cmd/qqBot/botMsg"
	"SamgeWxApi/cmd/qqBot/botUtil"
	"errors"
	"fmt"
)

// RunBot 运行机器人
func RunBot() {
	// 初始化
	if err := botUtil.BotEnvInit(); err != nil {
		panic(errors.New(fmt.Sprintf("机器人环境初始化失败：%s", err.Error())))
	}

	// 向监听器里添加函数
	botHandler.Listeners.OnPrivateMsg = append(botHandler.Listeners.OnPrivateMsg, botMsg.ParseMsgFriend)
	botHandler.Listeners.OnGroupMsg = append(botHandler.Listeners.OnGroupMsg, botMsg.ParseMsgGroup)

	bot := botHandler.CreateBot()
	botHandler.Run(bot)
}

func main() {
	RunBot()
}
