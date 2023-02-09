package botHandler

import (
	"SamgeWxApi/cmd/qqBot/botConfig"
	"SamgeWxApi/cmd/qqBot/botModel"
	"SamgeWxApi/cmd/utils/u_str"
	"os"
)

// CreateBot 创建机器人
func CreateBot() *botModel.Bot {
	bot := NewBot()
	bot.Config = botModel.Config{
		Loglvl:   LOGGER_LEVEL_INFO,
		Host:     u_str.FirstStr(botConfig.SocketHost, os.Getenv(botConfig.EnvKeySocketHost)),
		MasterQQ: u_str.Str2Int64(u_str.FirstStr(botConfig.MasterQQ, os.Getenv(botConfig.EnvKeyMasterQQ))),
		Path:     botConfig.SocketPath,
		MsgAwait: true,
		Token:    u_str.FirstStr(botConfig.SocketToken, os.Getenv(botConfig.EnvKeySocketToken)),
	}
	return bot
}
