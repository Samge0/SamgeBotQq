package botHandler

//! 整合大部分常用api于函数内

import (
	"SamgeWxApi/cmd/qqBot/botConfig"
	"SamgeWxApi/cmd/qqBot/botModel"
	"SamgeWxApi/cmd/utils/u_str"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var Connect *websocket.Conn // Websocket 连接接口

// ConnectWS 连接至websocket服务器
func ConnectWS(config *botModel.Config) (*websocket.Conn, *http.Response, error) {
	host := u_str.FirstStr(config.Host, os.Getenv(botConfig.EnvKeySocketHost))
	path := u_str.FirstStr(config.Path, os.Getenv(botConfig.SocketPath))
	token := u_str.FirstStr(config.Token, os.Getenv(botConfig.EnvKeyOpenAiToken))
	url := url.URL{Scheme: "ws", Host: host, Path: path}
	var dailer *websocket.Dialer
	header := map[string][]string{
		"Authorization": []string{token},
	}
	c, r, err := dailer.Dial(url.String(), header)
	return c, r, err
}

// sendwspack ws发包
func sendwspack(message string) error {
	if QqBot.Config.MsgAwait {
		rand.Seed(time.Now().Unix())
		time.Sleep(time.Duration(rand.Intn(5000)) * time.Millisecond)
	}
	err := Connect.WriteMessage(websocket.TextMessage, []byte(message))
	return err
}

// 发送API
func apiSend(apiType string, params string) (map[string]interface{}, error) {
	eventid := strconv.FormatInt(time.Now().UnixNano(), 10)
	ch := make(chan map[string]interface{})
	defer close(ch)

	chinfo := botModel.ShortEvent{Channel: &ch}

	botModel.ShortEvents[eventid] = chinfo

	Logger.Debug(fmt.Sprintf("[↑][EID:%s][Type:%s]S:%s", eventid, apiType, params))
	err := sendwspack(fmt.Sprintf(`{"action": "%s", "params": %s, "echo": "%s"}`, apiType, params, eventid))
	var receive map[string]interface{}
	if err == nil {
		select {
		case receive = <-ch:
			Logger.Debug(fmt.Sprintf("[↓][EID:%s][Type:%s]R:%s", eventid, apiType, receive))
		case <-time.After(5 * time.Second):
			Logger.Warning(fmt.Sprintf("[↓][EID:%s][Type:%s]Timeout", eventid, apiType))
			err = errors.New("timout in func apiSend")
		}
	} else {
		Logger.Warning(err.Error())
	}
	delete(botModel.ShortEvents, eventid)
	return receive, err
}

// onebot API: https://git.io/Jmy1B
// cqhttp API: https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#api-1

// SendPrivateMsg
// 发送私聊信息
// message - 消息内容 自动解析CQ码
// user_id - 对方QQ号
// return message_id error
func SendPrivateMsg(message string, user_id int64) (map[string]interface{}, error) {
	res, err := apiSend("send_private_msg", fmt.Sprintf(`{"user_id": %d, "message": "%s"}`, user_id, message))
	if err != nil {
		return nil, err
	}
	if res["status"].(string) == "ok" {
		Logger.Info(fmt.Sprintf("[↑][私聊][%d]: %s", user_id, message))
	} else {
		Logger.Warning(fmt.Sprintf("[↑][发送失败][私聊][%d]: %s", user_id, message))
	}
	return res, err
}

// SendGroupMsg
// 发送群聊消息
// message  - 要发送的内容
// group_id - 群号
// return message_id error
func SendGroupMsg(message string, group_id int64) (map[string]interface{}, error) {
	res, err := apiSend("send_group_msg", fmt.Sprintf(`{"group_id": %d, "message": "%s"}`, group_id, message))
	if err != nil {
		return nil, err
	}
	if res["status"].(string) == "ok" {
		Logger.Info(fmt.Sprintf("[↑][群聊][%d]: %s", group_id, message))
	} else {
		Logger.Warning(fmt.Sprintf("[↑][发送失败][群聊][%d]: %s", group_id, message))
	}
	return res, err
}

// SendMsg
// 发送消息
// msgtype - 消息类型 group/private
// message - 消息内容
// toid    - 群号/QQ号
// 本条API并不是 Onebot/CQhttp 原生API
// return message_id error
func SendMsg(msgtype string, message string, toid int64) (map[string]interface{}, error) {
	var err error
	var res map[string]interface{}
	switch msgtype {
	case "group":
		res, err = SendGroupMsg(message, toid)
	case "private":
		res, err = SendPrivateMsg(message, toid)
	default:
		return nil, errors.New("an error using function pichumod.SendMsg: msgtype should be group or private")
	}
	return res, err
}

// DeleteMsg
// 撤回消息
// message_id - 消息id 发出时的返回值
// return error
func DeleteMsg(message_id int32) error {
	_, err := apiSend("delete_msg", fmt.Sprintf(`{"message_id": %d}`, message_id))
	return err
}

// GetMsg
// 获取消息
// message_id - 获取消息
// return {time message_type message_id real_id sender message} error
func GetMsg(message_id int32) (map[string]interface{}, error) {
	res, err := apiSend("get_msg", fmt.Sprintf(`{"message_id": %d}`, message_id))
	return res, err
}

// GetForwardMsg
// 获取合并转发消息
// id - 合并转发 ID
// return message err
func GetForwardMsg(id string) (map[string]interface{}, error) {
	res, err := apiSend("get_forward_msg", fmt.Sprintf(`{"id":"%s"}`, id))
	return res, err
}

// SendLike
// 发送好友赞
// user_id - 对方QQ号
// times - 点赞次数(每个好友每天最多 10 次)
// return err
func SendLike(user_id int64, times int64) error {
	_, err := apiSend("send_like", fmt.Sprintf(`{"user_id": %d, "times": %d}`, user_id, times))
	return err
}

// SetGroupKick
// 群组踢人
// group_id - 群号
// user_id - 要踢的 QQ 号
// reject_add_request - 是否拒绝再次入群
// return err
func SetGroupKick(group_id int64, user_id int64, reject_add_request bool) error {
	_, err := apiSend("set_group_kick", fmt.Sprintf(`{"group_id": %d, "user_id": %d, "reject_add_request": %v}`, group_id, user_id, reject_add_request))
	return err
}

// SetGroupBan
// 群组单人禁言
// group_id - 群号
// user_id - 要禁言的QQ号
// duration - 禁言时长(s) 0表示取消禁言
// return err
func SetGroupBan(group_id int64, user_id int64, duration int64) error {
	_, err := apiSend("set_group_ban", fmt.Sprintf(`{"group_id": %d, "user_id": %d, "duration": "%d"}`, group_id, user_id, duration))
	return err
}

// SetGroupAnonymousBan
// 群组匿名用户禁言
// group_id - 群号
// anymous_flag - 匿名用户的 flag（需从群消息上报的数据中获得）
// duration - 禁言时长(s) 0表示取消禁言
// return err
func SetGroupAnonymousBan(group_id int64, anymous_flag string, duration int64) error {
	_, err := apiSend("set_group_ban", fmt.Sprintf(`{"group_id": %d, "anymous_flag": "%s", "duration": "%d"}`, group_id, anymous_flag, duration))
	return err
}

// SetGroupWholeBan
// 群全员禁言
// group_id 群号
// enable 是否禁言
// return err
func SetGroupWholeBan(group_id int64, enable bool) error {
	_, err := apiSend("set_group_kick", fmt.Sprintf(`{"group_id": %d, "enable": %v}`, group_id, enable))
	return err
}

// SetGroupAdmin
// 群组设置管理员(需要机器人为群主)
// group_id 群号
// user_id QQ号
// enable true 为设置，false 为取消
// return err
func SetGroupAdmin(group_id int64, user_id int64, enable bool) error {
	_, err := apiSend("set_group_admin", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "enable": %v}`, group_id, user_id, enable))
	return err
}

// SetGroupAnonymous
// 群组匿名
// group_id 群号
// enable 是否允许匿名聊天
// return err
func SetGroupAnonymous(group_id int64, enable bool) error {
	_, err := apiSend("set_group_anonymous", fmt.Sprintf(`{"group_id": %d, "enable": %v}`, group_id, enable))
	return err
}

// SetGroupCard
// 设置群名片
// group_id 群号
// user_id 成员QQ
// card 空字符串表示删除群名片
// return err
func SetGroupCard(group_id int64, user_id int64, card string) error {
	_, err := apiSend("set_group_card", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "card": "%s"}`, group_id, user_id, card))
	return err
}

// SetGroupName
// 设置群名
// group_id 群号
// group_name 新群名
// return err
func SetGroupName(group_id int64, group_name string) error {
	_, err := apiSend("set_group_name", fmt.Sprintf(`{"group_id": %d, "group_name": "%s"}`, group_id, group_name))
	return err
}

// SetGroupLeave
// 退群
// group_id 群号
// is_dismiss 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散
// return err
func SetGroupLeave(group_id int64, is_dismiss bool) error {
	_, err := apiSend("set_group_leave", fmt.Sprintf(`{"group_id": %d, "is_dismiss": %v}`, group_id, is_dismiss))
	return err
}

// SetGroupSpecialTitle
// 设置群组专属头衔
// group_id 群号
// user_id 成员QQ
// special_title 空字符串表示删除专属头衔
// return err
func SetGroupSpecialTitle(group_id int64, user_id int64, special_title string) error {
	_, err := apiSend("set_group_special_title", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "special_title": "%s"}`, group_id, user_id, special_title))
	return err
}

// SetFriendAddRequest
// 处理加好友请求
// flag 加好友请求的 flag（需从上报的数据中获得）
// approve 是否同意请求
// return err
func SetFriendAddRequest(flag string, approve bool) error {
	_, err := apiSend("set_friend_add_request", fmt.Sprintf(`{"flag": "%s", "approve": %v}`, flag, approve))
	return err
}

// SetGroupAddRequest
// 处理加群请求
// flag 加群请求的 flag
// approve 是否同意请求
// reason 拒绝理由(只有在拒绝时有效)
// return err
func SetGroupAddRequest(flag string, approve bool, reason string) error {
	_, err := apiSend("set_group_add_request", fmt.Sprintf(`{"flag": "%s", "sub_type": "add", "approve": %v, "reason": "%s"}`, flag, approve, reason))
	return err
}

// SetGroupInviteRequest
// 处理加群邀请
// flag 加群邀请的 flag
// approve 是否同意邀请
// reason 拒绝理由(只有在拒绝时有效)
// return err
func SetGroupInviteRequest(flag string, approve bool, reason string) error {
	_, err := apiSend("set_group_add_request", fmt.Sprintf(`{"flag": "%s", "sub_type": "invite", "approve": %v, "reason": "%s"}`, flag, approve, reason))
	return err
}

// GetLoginInfo
// 获取登录号信息
// return {user_id nickname} err
func GetLoginInfo() (map[string]interface{}, error) {
	res, err := apiSend("get_login_info", "")
	return res, err
}

// GetStrangerInfo
// 获取陌生人信息
// user_id 陌生人QQ
// no_cache 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func GetStrangerInfo(user_id int64, no_cache bool) (map[string]interface{}, error) {
	res, err := apiSend("get_stranger_info", fmt.Sprintf(`{"user_id": %d, "no_cache": %v}`, user_id, no_cache))
	return res, err
}

// GetFriendList
// 获取好友列表
// return [{user_id nickname}] err
func GetFriendList() (map[string]interface{}, error) {
	res, err := apiSend("get_friend_list", "")
	return res, err
}

// GetGroupInfo
// 获取群信息
// group_id 群号
// no_cache 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func GetGroupInfo(group_id int64, no_cache bool) (map[string]interface{}, error) {
	res, err := apiSend("get_group_info", fmt.Sprintf(`{"group_id": %d, "no_cache": %v}`, group_id, no_cache))
	return res, err
}

// GetGroupList
// 获取群列表
// return [{group_id name}] err
func GetGroupList() (map[string]interface{}, error) {
	res, err := apiSend("get_group_list", "")
	return res, err
}

// GetGroupMemberInfo
// 获取群成员信息
// group_id 群号
// user_id 成员QQ
// no_cache 是否不使用缓存（使用缓存可能更新不及时，但响应更快）
func GetGroupMemberInfo(group_id int64, user_id int64, no_cache bool) (map[string]interface{}, error) {
	res, err := apiSend("get_group_member_info", fmt.Sprintf(`{"group_id": %d, "user_id": %d, "no_cache": %v}`, group_id, user_id, no_cache))
	return res, err
}

// GetGroupMemberList
// 获取群成员列表
// group_id 群号
func GetGroupMemberList(group_id int64) (map[string]interface{}, error) {
	res, err := apiSend("get_group_member_list", fmt.Sprintf(`{"group_id": %d}`, group_id))
	return res, err
}

// GetGroupHonorInfo
// 获取群荣誉信息
// group_id 群号
// type 荣誉类型
func GetGroupHonorInfo(group_id int64, type_ int) (map[string]interface{}, error) {
	res, err := apiSend("get_group_honor_info", fmt.Sprintf(`{"group_id": %d, "type": %d}`, group_id, type_))
	return res, err
}

// 以下为 CQhttp 的API

// GetImage
// 获取图片信息
// file string
// get_image
// return size			int32		图片源文件大小
//
//	filename	string	图片文件原名
//	url				string	图片下载地址
//	err
func GetImage(file string) (map[string]interface{}, error) {
	res, err := apiSend("get_image", fmt.Sprintf(`{"file": "%s"}`, file))
	return res, err
}

// OCRImage
// 图片OCR
// image file - string
// return texts			TextDetection[]	OCR结果
//   - text				string	文本
//   - confidence	int32		置信度
//   - coordinates	vector2	坐标
//     language	string					语言
//     err
func OCRImage(image string) (map[string]interface{}, error) {
	res, err := apiSend("ocr_image", fmt.Sprintf(`{"image": "%s"}`, image))
	return res, err
}
