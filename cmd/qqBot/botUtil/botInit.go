package botUtil

import (
	"SamgeWxApi/cmd/qqBot/botConfig"
	"SamgeWxApi/cmd/utils/u_file"
)

// BotEnvInit 机器人环境初始化
func BotEnvInit() error {
	if err := InitCacheDir(); err != nil {
		return err
	}
	return nil
}

// InitCacheDir 初始化缓存目录
func InitCacheDir() error {
	if err := u_file.CreateMultiDir(botConfig.BotCacheDir); err != nil {
		return err
	}
	if err := u_file.CreateMultiDir(botConfig.BotLogDir); err != nil {
		return err
	}
	return nil
}
