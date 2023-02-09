package botHandler

import (
	botModel "SamgeWxApi/cmd/qqBot/botModel"
	"fmt"
)

// Listeners 监听器
var Listeners botModel.Listeners

// NewBot 返回一个Bot对象
func NewBot() *botModel.Bot {
	return &botModel.Bot{}
}

// MsgParse 处理接收到的消息
func MsgParse(receive map[string]interface{}) {

	switch receive["post_type"] {
	// 消息事件
	case "message":
		switch receive["message_type"] {
		// 私聊信息
		case "private":
			var eventinfo botModel.MessagePrivate = parsePrivate(receive)
			Logger.Info(fmt.Sprintf("[↓][私聊][%s(%d)]: %s", eventinfo.Sender.Nickname, eventinfo.Sender.UserID, eventinfo.Message))
			for _, function := range Listeners.OnPrivateMsg {
				function(eventinfo)
			}

		// 群聊信息
		case "group":
			var eventinfo botModel.MessageGroup = parseGroup(receive)
			Logger.Info(fmt.Sprintf("[↓][群聊(%d)][%s(%d)]: %s", eventinfo.GroupID, eventinfo.Msg.Sender.Nickname, eventinfo.Msg.Sender.UserID, eventinfo.Msg.Message))
			for _, function := range Listeners.OnGroupMsg {
				go function(eventinfo)
			}

		default:
			Logger.Warning(fmt.Sprintf("Cannot Parse 'message' event -> %s", receive))
		}

		// 通知事件
	case "notice":
		switch receive["notice_type"] {
		// 群文件上传
		case "group_upload":
			var eventinfo botModel.GroupUpload = parseGroupupload(receive)
			Logger.Info(fmt.Sprintf("[N][群文件(%d)][%d]: %s", eventinfo.GroupId, eventinfo.UserId, eventinfo.File.Name))
			for _, function := range Listeners.OnGroupUpload {
				go function(eventinfo)
			}

			// 群管理员变动
		case "group_admin":
			var eventinfo botModel.GroupAdmin = parseGroupadmin(receive)
			var x string
			if eventinfo.SubType == "set" {
				x = "+"
			} else {
				x = "-"
			}
			Logger.Info(fmt.Sprintf("[N][群(%d)管理][%s %d]", eventinfo.GroupId, x, eventinfo.UserId))
			for _, function := range Listeners.OnGroupAdmin {
				go function(eventinfo)
			}

			// 群成员减少
		case "group_decrease":
			var eventinfo botModel.GroupDecrease = parseGroupdecrease(receive)
			Logger.Info(fmt.Sprintf("[N][成员退群(%d)][%d] Type: %s", eventinfo.GroupId, eventinfo.UserId, eventinfo.SubType))
			for _, function := range Listeners.OnGroupDecrease {
				go function(eventinfo)
			}

			// 群成员增加
		case "group_increase":
			var eventinfo botModel.GroupIncrease = parseGroupincrease(receive)
			Logger.Info(fmt.Sprintf("[N][成员入群(%d)][%d -> %d] Type: %s", eventinfo.GroupId, eventinfo.OperatorId, eventinfo.UserId, eventinfo.SubType))
			for _, function := range Listeners.OnGroupIncrease {
				go function(eventinfo)
			}

			// 群禁言
		case "group_ban":
			var eventinfo botModel.GroupBan = parseGroupban(receive)
			Logger.Info(fmt.Sprintf("[N][群聊(%d)] %d 禁言/解禁了 %d for %ds", eventinfo.GroupId, eventinfo.OperatorId, eventinfo.UserId, eventinfo.Duration))
			for _, function := range Listeners.OnGroupBan {
				go function(eventinfo)
			}

			// 好友添加
		case "friend_add":
			var eventinfo botModel.FriendAdd = parseFriendAdd(receive)
			Logger.Info(fmt.Sprintf("[N][成功添加好友]%d", eventinfo.UserId))
			for _, function := range Listeners.OnFriendAdd {
				go function(eventinfo)
			}

			// 群消息撤回
		case "group_recall":
			var eventinfo botModel.GroupRecall = parseGrouprecall(receive)
			Logger.Info(fmt.Sprintf("[N][群聊(%d)][%d] 撤回了消息(id): %d", eventinfo.GroupId, eventinfo.UserId, eventinfo.MessageId))
			for _, function := range Listeners.OnGroupRecall {
				go function(eventinfo)
			}

			// 好友消息撤回
		case "friend_recall":
			var eventinfo botModel.FriendRecall = parseFriendrecall(receive)
			Logger.Info(fmt.Sprintf("[N][私聊][%d] 撤回了消息(id): %d", eventinfo.UserId, eventinfo.MessageId))
			for _, function := range Listeners.OnFriendRecall {
				go function(eventinfo)
			}

			// 群内戳一戳 群红包运气王 群成员荣誉变更
		case "notify":
			var eventinfo botModel.Notify = parseNotify(receive)
			Logger.Info(fmt.Sprintf("[N][Notify][Group:%d] %d -> %s", eventinfo.GroupId, eventinfo.UserId, eventinfo.SubType))
			for _, function := range Listeners.OnNotify {
				go function(eventinfo)
			}

		default:
			Logger.Warning(fmt.Sprintf("Cannot Parse 'notice' event -> %s", receive))
		}

		// 请求事件
	case "request":
		switch receive["request_type"] {
		// 添加好友申请
		case "friend":
			var eventinfo botModel.FriendRequest = parseFriendrequest(receive)
			Logger.Info(fmt.Sprintf("[↓][好友申请] %d 申请加你为好友 -> %s", eventinfo.UserId, eventinfo.Comment))
			for _, function := range Listeners.OnFriendRequest {
				go function(eventinfo)
			}

			// 加群邀请
		case "group":
			// SetGroupInviteRequest(receive["flag"].(string), true, "") // 自动同意加群
			var eventinfo botModel.GroupRequest = parseGrouprequest(receive)
			Logger.Info(fmt.Sprintf("[↓][加群/邀请] %d %s -> %d(验证信息: %s)", eventinfo.UserId, eventinfo.SubType, eventinfo.GroupId, eventinfo.Comment))
			for _, function := range Listeners.OnGroupRequest {
				go function(eventinfo)
			}

		default:
			Logger.Warning(fmt.Sprintf("Cannot Parse 'request' event -> %s", receive))
		}
		// 元事件
	case "meta_event":
		switch receive["meta_event_type"] {
		// 生命周期
		case "lifecycle":
			var eventinfo botModel.MetaLifecycle = parseMetalifecycle(receive)
			Logger.Debug(fmt.Sprintf("[↓][Lifecycle][%d] Type: %s", eventinfo.SelfId, eventinfo.SubType))
			for _, function := range Listeners.OnMetaLifecycle {
				go function(eventinfo)
			}

			// 心跳包
		case "heartbeat":
			var eventinfo botModel.MetaHeartbeat = parseMetaheartbeat(receive)
			Logger.Debug(fmt.Sprintf("[↓][Heartbeat][%d] Type: %s", eventinfo.SelfId, eventinfo.Status))
			for _, function := range Listeners.OnMetaHeartbeat {
				go function(eventinfo)
			}

			// Logger.Debug("Received a heartbeat package.")
		default:
			Logger.Warning(fmt.Sprintf("Cannot Parse 'meta_event' event -> %s", receive))
		}
	default:
		// 短事件回调
		if _, ok := receive["echo"]; ok {
			if _, ok := botModel.ShortEvents[receive["echo"].(string)]; ok {
				*botModel.ShortEvents[receive["echo"].(string)].Channel <- receive
			}
		} else {
			Logger.Warning(fmt.Sprintf("Got Error Package -> %s", receive))
		}
	}
}

func parsePrivate(r map[string]interface{}) botModel.MessagePrivate {
	e := botModel.MessagePrivate{
		SelfID:     int64(r["self_id"].(float64)),
		SubType:    r["sub_type"].(string),
		MessageID:  int64(r["message_id"].(float64)),
		UserID:     int64(r["user_id"].(float64)),
		Message:    r["message"].(string),
		RawMessage: r["raw_message"].(string),
		Sender: botModel.Sender{
			UserID:   int64(r["sender"].(map[string]interface{})["user_id"].(float64)),
			Nickname: r["sender"].(map[string]interface{})["nickname"].(string),
			Sex:      r["sender"].(map[string]interface{})["sex"].(string),
			Age:      int64(r["sender"].(map[string]interface{})["age"].(float64)),
		}}
	return e
}

// parseGroup 处理群组消息
func parseGroup(r map[string]interface{}) botModel.MessageGroup {
	e := botModel.MessageGroup{
		Msg: &botModel.MessagePrivate{
			SelfID:     int64(r["self_id"].(float64)),
			SubType:    r["sub_type"].(string),
			MessageID:  int64(r["message_id"].(float64)),
			UserID:     int64(r["user_id"].(float64)),
			Message:    r["message"].(string),
			RawMessage: r["raw_message"].(string),
		},
		GroupID: int64(r["group_id"].(float64)),
	}

	switch e.Msg.SubType {
	case "normal":
		e.Msg.Sender = botModel.Sender{
			UserID:   int64(r["sender"].(map[string]interface{})["user_id"].(float64)),
			Nickname: r["sender"].(map[string]interface{})["nickname"].(string),
			Card:     r["sender"].(map[string]interface{})["card"].(string),
			Sex:      r["sender"].(map[string]interface{})["sex"].(string),
			Age:      int64(r["sender"].(map[string]interface{})["age"].(float64)),
			Area:     r["sender"].(map[string]interface{})["area"].(string),
			Level:    r["sender"].(map[string]interface{})["level"].(string),
			Role:     r["sender"].(map[string]interface{})["role"].(string),
			Title:    r["sender"].(map[string]interface{})["title"].(string)}
	case "anoymous":
		e.Anonymous = struct {
			Id   int64
			Name string
			Flag string
		}{
			Id:   int64(r["anonymous"].(map[string]interface{})["id"].(float64)),
			Name: r["anonymous"].(map[string]interface{})["name"].(string),
			Flag: r["anonymous"].(map[string]interface{})["flag"].(string)}
	}
	return e
}
func parseGroupupload(r map[string]interface{}) botModel.GroupUpload {
	e := botModel.GroupUpload{
		Time:    int64(r["time"].(float64)),
		SelfId:  int64(r["self_id"].(float64)),
		GroupId: int64(r["group_id"].(float64)),
		UserId:  int64(r["user_id"].(float64)),
		File: struct {
			Id    string
			Name  string
			Size  int64
			Busid int64
		}{
			Id:    r["file"].(map[string]interface{})["id"].(string),
			Name:  r["file"].(map[string]interface{})["name"].(string),
			Size:  int64(r["file"].(map[string]interface{})["size"].(float64)),
			Busid: int64(r["file"].(map[string]interface{})["busid"].(float64)),
		}}
	return e
}
func parseGroupadmin(r map[string]interface{}) botModel.GroupAdmin {
	e := botModel.GroupAdmin{
		Time:    int64(r["time"].(float64)),
		SelfId:  int64(r["self_id"].(float64)),
		SubType: r["sub_type"].(string),
		GroupId: int64(r["group_id"].(float64)),
		UserId:  int64(r["user_id"].(float64)),
	}
	return e
}
func parseGroupdecrease(r map[string]interface{}) botModel.GroupDecrease {
	e := botModel.GroupDecrease{
		Time:       int64(r["time"].(float64)),
		SelfId:     int64(r["self_id"].(float64)),
		SubType:    r["sub_type"].(string),
		GroupId:    int64(r["group_id"].(float64)),
		OperatorId: int64(r["operator_id"].(float64)),
		UserId:     int64(r["user_id"].(float64)),
	}
	return e
}
func parseGroupincrease(r map[string]interface{}) botModel.GroupIncrease {
	e := botModel.GroupIncrease{
		Time:       int64(r["time"].(float64)),
		SelfId:     int64(r["self_id"].(float64)),
		SubType:    r["sub_type"].(string),
		GroupId:    int64(r["group_id"].(float64)),
		OperatorId: int64(r["operator_id"].(float64)),
		UserId:     int64(r["user_id"].(float64)),
	}
	return e
}
func parseGroupban(r map[string]interface{}) botModel.GroupBan {
	e := botModel.GroupBan{
		Time:       int64(r["time"].(float64)),
		SelfId:     int64(r["self_id"].(float64)),
		SubType:    r["sub_type"].(string),
		GroupId:    int64(r["group_id"].(float64)),
		OperatorId: int64(r["operator_id"].(float64)),
		UserId:     int64(r["user_id"].(float64)),
		Duration:   int64(r["duration"].(float64)),
	}
	return e
}
func parseFriendAdd(r map[string]interface{}) botModel.FriendAdd {
	e := botModel.FriendAdd{
		Time:   int64(r["time"].(float64)),
		SelfId: int64(r["self_id"].(float64)),
		UserId: int64(r["user_id"].(float64)),
	}
	return e
}

func parseGrouprecall(r map[string]interface{}) botModel.GroupRecall {
	e := botModel.GroupRecall{
		Time:       int64(r["time"].(float64)),
		SelfId:     int64(r["self_id"].(float64)),
		GroupId:    int64(r["group_id"].(float64)),
		UserId:     int64(r["user_id"].(float64)),
		OperatorId: int64(r["operator_id"].(float64)),
		MessageId:  int64(r["message_id"].(float64)),
	}
	return e
}
func parseFriendrecall(r map[string]interface{}) botModel.FriendRecall {
	e := botModel.FriendRecall{
		Time:      int64(r["time"].(float64)),
		SelfId:    int64(r["self_id"].(float64)),
		UserId:    int64(r["user_id"].(float64)),
		MessageId: int64(r["message_id"].(float64)),
	}
	return e
}
func parseNotify(r map[string]interface{}) botModel.Notify {
	e := botModel.Notify{
		Time:    int64(r["time"].(float64)),
		SelfId:  int64(r["self_id"].(float64)),
		SubType: r["sub_type"].(string),
		GroupId: int64(r["group_id"].(float64)),
		UserId:  int64(r["user_id"].(float64)),
	}
	if e.SubType == "honor" {
		e.Honor_type = r["honor_type"].(string)
	} else {
		e.TargetId = int64(r["target_id"].(float64))
	}
	return e
}
func parseFriendrequest(r map[string]interface{}) botModel.FriendRequest {
	e := botModel.FriendRequest{
		Time:    int64(r["time"].(float64)),
		SelfId:  int64(r["self_id"].(float64)),
		UserId:  int64(r["user_id"].(float64)),
		Comment: r["comment"].(string),
		Flag:    r["flag"].(string),
	}
	return e
}
func parseGrouprequest(r map[string]interface{}) botModel.GroupRequest {
	e := botModel.GroupRequest{
		Time:    int64(r["time"].(float64)),
		SelfId:  int64(r["self_id"].(float64)),
		SubType: r["sub_type"].(string),
		GroupId: int64(r["group_id"].(float64)),
		UserId:  int64(r["user_id"].(float64)),
		Comment: r["comment"].(string),
		Flag:    r["flag"].(string),
	}
	return e
}

func parseMetalifecycle(r map[string]interface{}) botModel.MetaLifecycle {
	e := botModel.MetaLifecycle{
		Time:    int64(r["time"].(float64)),
		SelfId:  int64(r["self_id"].(float64)),
		SubType: r["sub_type"].(string),
	}
	return e
}
func parseMetaheartbeat(r map[string]interface{}) botModel.MetaHeartbeat {
	e := botModel.MetaHeartbeat{
		Time:     int64(r["time"].(float64)),
		SelfId:   int64(r["self_id"].(float64)),
		Status:   r["status"],
		Interval: int64(r["interval"].(float64)),
	}
	return e
}
