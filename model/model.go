package model

import "DouyinLive/database"

type Comment struct {
	LiveId  int    `json:"live_id"`
	RoomId  int    `json:"room_id"`
	UserId  int    `json:"user_id"`
	Content string `json:"content"`
}

func InsertComments(liveId, roomId, userId int, content string) {
	comment := Comment{
		LiveId:  liveId,
		RoomId:  roomId,
		UserId:  userId,
		Content: content,
	}
	database.DB.Table("comments").Create(&comment)
}
