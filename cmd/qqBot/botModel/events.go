package botModel

import (
	"strconv"
	"time"
)

// ShortEvent 短事件
type ShortEvent struct {
	Channel *chan map[string]interface{}
}

// LongEvent 长事件
type LongEvent struct {
	UserID   int64
	GroupID  int64
	Channel  *chan string
	EventKey string
	EventID  string
}

var ShortEvents = make(map[string]ShortEvent) // 短事件容器
var LongEvents = make(map[string]LongEvent)   // 长事件容器

// NewEvent 新建长事件
func NewEvent(userid int64, groupid int64, key string) LongEvent {
	var eventId string = strconv.FormatInt(time.Now().UnixNano(), 10) // 事件ID 以时间戳定义
	ch := make(chan string)
	event := LongEvent{Channel: &ch, UserID: userid, GroupID: groupid, EventKey: key, EventID: eventId}
	LongEvents[eventId] = event
	return event
}

// Close 关闭事件
func (event LongEvent) Close() {
	close(*event.Channel)
	delete(LongEvents, event.EventID)
}
