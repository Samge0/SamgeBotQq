package botHandler

import (
	"SamgeWxApi/cmd/qqBot/botConfig"
	"SamgeWxApi/cmd/qqBot/botModel"
	"SamgeWxApi/cmd/utils/u_openai"
	"SamgeWxApi/cmd/utils/u_str"
	"fmt"
	"github.com/eatmoreapple/openwechat"
	"os"
	"strings"
)

// ReplyWithOpenAi 使用openai的api进行回复
func ReplyWithOpenAi(msg *botModel.MessageParam, question string, qType string) {
	openAiToken := u_str.FirstStr(botConfig.OpenAiToken, os.Getenv(botConfig.EnvKeyOpenAiToken))
	answer, err := u_openai.GetChatResponseWithToken(question, 1000, openAiToken)
	if err != nil {
		SaveErrorLog(err, "ReplyWithOpenAi："+qType)
		answer = "啊？你说什么？风太大我没听清，请再说一遍"
	}
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), question, answer, qType)
}

// CheckStartTagAndReply 检查内容的起始标签，如果符合则进行回复
// 普通文本：hi
// 表情：[CQ:face,id=212]
// 群@：[CQ:at,qq=1039476587] hi
func CheckStartTagAndReply(msg *botModel.MessageParam, qType string) {
	question := msg.Msg.Message

	result := checkMasterOrder(question, msg.Msg.Sender.UserID)
	if result != "" {
		ReplyText(msg, result)
		return
	}

	if IsStopChat {
		return
	}

	fmt.Println(fmt.Sprintf("收到消息：%s", question))
	openAiToken := u_str.FirstStr(botConfig.OpenAiToken, os.Getenv(botConfig.EnvKeyOpenAiToken))
	answer, err := u_openai.GetChatResponseWithToken(question, 1000, openAiToken)
	if err != nil {
		SaveErrorLog(err, "ReplyWithOpenAi："+qType)
		fmt.Println(fmt.Sprintf("ReplyWithOpenAi：%s", err.Error()))
		answer = "啊？你说什么？风太大我没听清，请再说一遍"
	}
	ReplyText(msg, answer)

	/*if msg.IsTickledMe() {
		ParseMsgOnTickled(msg, fmt.Sprintf("%s 拍一拍", qType))
	} else if msg.IsText() { // 文本
		ParseMsgOnText(msg, qType)
	} else if msg.IsPicture() { // 图片
		ParseMsgOnImage(msg, qType)
	} else if msg.IsVoice() { // 语音
		ParseMsgOnVoice(msg, qType)
	} else if msg.IsCard() { // 卡片
		ParseMsgOnCard(msg, qType)
	} else if msg.IsVideo() { // 视频
		ParseMsgOnVideo(msg, qType)
	} else if msg.IsEmoticon() { // 表情包
		ParseMsgOnEmoticon(msg, qType)
	} else if msg.IsRealtimeLocation() { // 实时位置共享
		ParseMsgOnRealtimeLocation(msg, qType)
	} else if msg.IsLocation() { // 位置
		ParseMsgOnLocation(msg, qType)
	} else if msg.IsTransferAccounts() { // 微信转账
		ParseMsgOnTransferAccounts(msg, qType)
	} else if msg.IsSendRedPacket() { // 微信红包-发出
		ParseMsgOnSendRedPacket(msg, qType)
	} else if msg.IsReceiveRedPacket() { // 微信红包-收到
		ParseMsgOnReceiveRedPacket(msg, qType)
	} else if msg.IsRenameGroup() { // 群组重命名
		ParseMsgOnRenameGroup(msg, qType)
	} else if msg.IsArticle() { // 文章
		ParseMsgOnArticle(msg, qType)
	} else if msg.IsVoipInvite() { // 语音/视频邀请
		ParseMsgOnVoipInvite(msg, qType)
	} else if msg.IsMedia() { // Media(多媒体消息，包括但不限于APP分享、文件分享
		ParseMsgOnMedia(msg, qType)
	}*/

}

// checkMasterOrder 检查管理者的命令
func checkMasterOrder(question string, userId int64) string {
	masterQQ := u_str.FirstStr(botConfig.MasterQQ, os.Getenv(botConfig.EnvKeyMasterQQ))
	if !strings.Contains(masterQQ, fmt.Sprintf("%d", userId)) {
		return "哎呀呀，您的指令被三体人吃掉了"
	}
	defaultReply := fmt.Sprintf("收到！已%s", question)
	switch question {
	case "退下":
		IsStopChat = true
		return defaultReply
	case "恢复":
		IsStopChat = false
		return defaultReply
	default:
		return ""
	}
}

// ParseMsgOnTickled 处理【拍一拍】类型的消息
func ParseMsgOnTickled(msg *botModel.MessageParam, qType string) {
	answer := "再拍我我就把你吃了"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), msg.Msg.Message, answer, qType)
}

// ParseMsgOnText 处理【OnText】类型的消息
func ParseMsgOnText(msg *botModel.MessageParam, qType string) {
	tagHead := "生成头像 "
	tagT2I := "生成图片 "

	var question string
	question = msg.Msg.Message
	if strings.HasPrefix(question, tagHead) {
		question = strings.Replace(question, tagHead, "", 1)
		answer := fmt.Sprintf("%s 服务正在开发中", tagHead)
		ReplyText(msg, answer)
		SaveChatLog(msg.Msg.Sender.String(), question, answer, tagHead)
	} else if strings.HasPrefix(question, tagT2I) {
		question = strings.Replace(question, tagT2I, "", 1)
		answer := fmt.Sprintf("%s 服务正在开发中", tagT2I)
		ReplyText(msg, answer)
		SaveChatLog(msg.Msg.Sender.String(), question, answer, tagT2I)
	} else {
		ReplyWithOpenAi(msg, question, qType)
	}
}

// ParseMsgOnImage 处理【OnImage】类型的消息
func ParseMsgOnImage(msg *botModel.MessageParam, qType string) {
	answer := "这是什么图片"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 图片", qType))

	//fileName := u_date.GetCurrentDateStr(u_date.DateFormat.Flow)
	//msg.SaveFileToLocal(fmt.Sprintf("%s/%s_%s.jpg", botConfig.BotCacheDir, sender.NickName, u_str.FirstStr(msg.FileName, fileName)))
}

// ParseMsgOnVoice 处理【OnVoice】类型的消息
func ParseMsgOnVoice(msg *botModel.MessageParam, qType string) {
	answer := "不方便接听语音信息，还是发文字吧"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 语音", qType))

	//fileName := u_date.GetCurrentDateStr(u_date.DateFormat.Flow)
	//msg.SaveFileToLocal(fmt.Sprintf("%s/%s_%s.amr", botConfig.BotCacheDir, sender.NickName, u_str.FirstStr(msg.FileName, fileName)))
}

// ParseMsgOnCard 处理【OnCard】类型的消息
func ParseMsgOnCard(msg *botModel.MessageParam, qType string) {
	answer := "这是什么好玩的小卡片"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 卡片", qType))
}

// ParseMsgOnMedia 处理【Media(多媒体消息，包括但不限于APP分享、文件分享)的处理函数】类型的消息
func ParseMsgOnMedia(msg *botModel.MessageParam, qType string) {
	answer := "你这发的是个啥子哦"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 多媒体消息等", qType))

	//fileName := u_date.GetCurrentDateStr(u_date.DateFormat.Flow)
	//msg.SaveFileToLocal(fmt.Sprintf("%s/%s_%s.file", botConfig.BotCacheDir, sender.NickName, u_str.FirstStr(msg.FileName, fileName)))
}

// ParseMsgOnVideo 处理【视频】类型的消息
func ParseMsgOnVideo(msg *botModel.MessageParam, qType string) {
	answer := "这是什么视频，好看的话多发点来look look"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 视频", qType))

	//fileName := u_date.GetCurrentDateStr(u_date.DateFormat.Flow)
	//msg.SaveFileToLocal(fmt.Sprintf("%s/%s_%s.mp4", botConfig.BotCacheDir, sender.NickName, u_str.FirstStr(msg.FileName, fileName)))
}

// ParseMsgOnEmoticon 处理【表情】类型的消息
func ParseMsgOnEmoticon(msg *botModel.MessageParam, qType string) {
	answer := fmt.Sprintf("%s%s", openwechat.Emoji.Awesome, openwechat.Emoji.Doge)
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 表情", qType))

	//fileName := u_date.GetCurrentDateStr(u_date.DateFormat.Flow)
	//msg.SaveFileToLocal(fmt.Sprintf("%s/%s_%s.gif", botConfig.BotCacheDir, sender.NickName, u_str.FirstStr(msg.FileName, fileName)))
}

// ParseMsgOnRealtimeLocation 处理【实时位置】类型的消息
func ParseMsgOnRealtimeLocation(msg *botModel.MessageParam, qType string) {
	answer := "你现在在哪个位置？我点不开看不到"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 实时位置", qType))
}

// ParseMsgOnLocation 处理【位置】类型的消息
func ParseMsgOnLocation(msg *botModel.MessageParam, qType string) {
	answer := "这是哪里？你又到哪里鬼混去了"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 位置", qType))
}

// ParseMsgOnTransferAccounts 处理【微信转账】类型的消息
func ParseMsgOnTransferAccounts(msg *botModel.MessageParam, qType string) {
	answer := "多谢老板，祝老板身体健康发大财~"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 微信转账", qType))
}

// ParseMsgOnSendRedPacket 处理【微信红包-发出】类型的消息
func ParseMsgOnSendRedPacket(msg *botModel.MessageParam, qType string) {
}

// ParseMsgOnReceiveRedPacket 处理【微信红包-收到】类型的消息
func ParseMsgOnReceiveRedPacket(msg *botModel.MessageParam, qType string) {
	answer := "多谢老板的大红包，好事成双，再发一个吧~"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 微信红包-收到", qType))
}

// ParseMsgOnRenameGroup 处理【群组重命名】类型的消息
func ParseMsgOnRenameGroup(msg *botModel.MessageParam, qType string) {
	answer := "群名变来变去的累不累哦？"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 群组重命名", qType))
}

// ParseMsgOnArticle 处理【文章消息】类型的消息
func ParseMsgOnArticle(msg *botModel.MessageParam, qType string) {
	answer := "这是什么绝世好文"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 文章", qType))
}

// ParseMsgOnVoipInvite 处理【语音或视频通话邀请】类型的消息
func ParseMsgOnVoipInvite(msg *botModel.MessageParam, qType string) {
	answer := "我现在在忙，不方便接听"
	ReplyText(msg, answer)
	SaveChatLog(msg.Msg.Sender.String(), "", answer, fmt.Sprintf("%s 通话邀请", qType))
}

// ReplyText 回复文本，如果是群聊，则@对方
func ReplyText(msg *botModel.MessageParam, message string) {
	var result map[string]interface{}
	var err error
	if msg.Msg.IsGroup() {
		result, err = SendGroupMsg(fmt.Sprintf("@%s %s", msg.Msg.Sender.Nickname, message), msg.GroupID)
	} else {
		result, err = SendPrivateMsg(message, msg.Msg.Sender.UserID)
	}
	if err != nil {
		fmt.Println(fmt.Sprintf("发送失败：%s", err.Error()))
	} else {
		fmt.Println(fmt.Sprintf("发送成功：%s", result))
	}
}
