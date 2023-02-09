package botMsg

import (
	"SamgeWxApi/cmd/qqBot/botConfig"
	"SamgeWxApi/cmd/qqBot/botHandler"
	"SamgeWxApi/cmd/qqBot/botModel"
	"SamgeWxApi/cmd/utils/u_str"
	"fmt"
	"os"
	"strings"
)

// ParseMsgFriend 处理朋友私聊消息
func ParseMsgFriend(eventInfo botModel.MessagePrivate) {
	userId := eventInfo.Sender.UserID
	userIdStr := fmt.Sprintf("%d", userId)
	ids := u_str.FirstStr(botConfig.FriendIds, os.Getenv(botConfig.EnvKeyFriendIds))
	needParseMsg := ids == "" || strings.Contains(ids, userIdStr)
	if needParseMsg {
		qType := fmt.Sprintf("[好友]%s|%s|%d", eventInfo.Sender.Nickname, eventInfo.Sender.Sex, userId)
		botHandler.CheckStartTagAndReply(eventInfo.AsMessageParam(), qType)
	}
}
