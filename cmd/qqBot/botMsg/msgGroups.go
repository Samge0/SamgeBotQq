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

// ParseMsgGroup 处理群组消息
// 群@：[CQ:at,qq=1039476587] hi
func ParseMsgGroup(eventInfo botModel.MessageGroup) {
	groupId := eventInfo.GroupID
	groupIdStr := fmt.Sprintf("%d", groupId)
	ids := u_str.FirstStr(botConfig.GroupIds, os.Getenv(botConfig.EnvKeyGroupIds))
	question := eventInfo.Msg.Message

	atMeTag := fmt.Sprintf("[CQ:at,qq=%d] ", eventInfo.Msg.SelfID)
	isAtMe := strings.Contains(question, atMeTag)
	needParseMsg := isAtMe && (ids == "" || strings.Contains(ids, groupIdStr))

	if needParseMsg {
		eventInfo.Msg.Message = strings.Replace(eventInfo.Msg.Message, atMeTag, "", 1)
		eventInfo.Msg.SubType = "group"
		qType := fmt.Sprintf("[群组]%s|%s|%d|%s", eventInfo.Msg.Sender.Nickname, eventInfo.Msg.Sender.Sex, eventInfo.Msg.Sender.UserID, groupIdStr)
		botHandler.CheckStartTagAndReply(eventInfo.AsMessageParam(), qType)
	}
}
