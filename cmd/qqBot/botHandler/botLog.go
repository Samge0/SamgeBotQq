package botHandler

import (
	"SamgeWxApi/cmd/qqBot/botConfig"
	"SamgeWxApi/cmd/utils/u_bot"
	"SamgeWxApi/cmd/utils/u_file"
	"fmt"
	goLogger "github.com/phachon/go-logger"
)

const (
	LOGGER_LEVEL_EMERGENCY = iota // 系统级紧急，比如磁盘出错，内存异常，网络不可用等
	LOGGER_LEVEL_ALERT            // 系统级警告，比如数据库访问异常，配置文件出错等
	LOGGER_LEVEL_CRITICAL         // 系统级危险，比如权限出错，访问异常等
	LOGGER_LEVEL_ERROR            // 用户级错误
	LOGGER_LEVEL_WARNING          // 用户级警告
	LOGGER_LEVEL_NOTICE           // 用户级重要
	LOGGER_LEVEL_INFO             // 用户级提示
	LOGGER_LEVEL_DEBUG            // 用户级调试
)

// Logger 日志接口
var Logger *goLogger.Logger

// initLogPath 初始化日志路径
func initLogPath() error {
	if err := u_file.CreateMultiDir(botConfig.BotLogDir); err != nil {
		return err
	}
	return nil
}

// InitLogger 初始化日志
func InitLogger(lvl int) error {
	if err := initLogPath(); err != nil {
		return err
	}

	Logger = goLogger.NewLogger()
	_ = Logger.Detach("console") // 禁用默认控制台日志

	consoleConfig := &goLogger.ConsoleConfig{
		Color:      true,
		JsonFormat: false,
		Format:     "%timestamp_format% [%level_string%] %body%",
	}
	_ = Logger.Attach("console", lvl, consoleConfig) // 加载指定规则的控制台日志

	fileConfig := &goLogger.FileConfig{
		Filename:   fmt.Sprintf("%s/qqBot.log", botConfig.BotLogDir),
		MaxSize:    1024 * 1024,
		MaxLine:    10000,
		DateSlice:  "d",
		JsonFormat: false,
		Format:     "%millisecond_format% [%level_string%] [%file%:%line%] %body%",
	}
	_ = Logger.Attach("file", lvl, fileConfig) // 配置日志保存文件

	Logger.SetAsync()
	return nil
}

// SaveErrorLog 保存错误日志
func SaveErrorLog(content any, contentType string) {
	u_bot.SaveErrorLog(botConfig.BotLogDir, content, contentType)
}

// SaveChatLog 保存聊天日志
func SaveChatLog(sengInfo string, question string, answer string, qType string) {
	u_bot.SaveChatLog(botConfig.BotLogDir, sengInfo, question, answer, qType)
}
