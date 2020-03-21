package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type Message struct {
	ID        int64     `db:"id"`
	ChannelID int64     `db:"channel_id"`
	UserID    int64     `db:"user_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
}

func postMessage(c echo.Context) error {
	user, err := ensureLogin(c)
	if user == nil {
		return err
	}

	message := c.FormValue("message")
	if message == "" {
		return echo.ErrForbidden
	}

	var chanID int64
	if x, err := strconv.Atoi(c.FormValue("channel_id")); err != nil {
		return echo.ErrForbidden
	} else {
		chanID = int64(x)
	}

	if _, err := addMessage(chanID, user.ID, message); err != nil {
		return err
	}

	return c.NoContent(204)
}

func queryMessages(chanID, lastID int64) ([]Message, error) {
	msgs := []Message{}
	err := db.Select(&msgs, "SELECT * FROM message WHERE id > ? AND channel_id = ? ORDER BY id DESC LIMIT 100",
		lastID, chanID)
	return msgs, err
}

type MsgWithUsr struct {
	User User
	Msg  Message
}

func queryMsgWithUsrs(chanID, lastID int64) ([]MsgWithUsr, error) {
	msgWithUsrs := []MsgWithUsr{}

	query := `SELECT u.name, u.display_name, u.avatar_icon, m.id, m.created_at, m.content FROM user u INNER JOIN message m ON u.id = m.user_id WHERE m.id > ? AND m.channel_id = ? ORDER BY m.id DESC LIMIT 100`
	rows, err := db.Query(query, lastID, chanID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		u := User{}
		m := Message{}

		rows.Scan(&u.Name, &u.DisplayName, &u.AvatarIcon, &m.ID, &m.CreatedAt, &m.Content)
		if err != nil {
			return nil, err
		}

		msgWithUsr := MsgWithUsr{User: u, Msg: m}

		msgWithUsrs = append(msgWithUsrs, msgWithUsr)
	}

	return msgWithUsrs, nil
}

func pureJsonifyMessage(mwu MsgWithUsr) map[string]interface{} {
	r := make(map[string]interface{})
	r["id"] = mwu.Msg.ID
	r["user"] = mwu.User
	r["date"] = mwu.Msg.CreatedAt.Format("2006/01/02 15:04:05")
	r["content"] = mwu.Msg.Content
	return r
}

func jsonifyMessage(m Message) (map[string]interface{}, error) {
	u := User{}

	err := db.Get(&u, "SELECT name, display_name, avatar_icon FROM user WHERE id = ?",
		m.UserID)
	if err != nil {
		return nil, err
	}

	r := make(map[string]interface{})
	r["id"] = m.ID
	r["user"] = u
	r["date"] = m.CreatedAt.Format("2006/01/02 15:04:05")
	r["content"] = m.Content
	return r, nil
}

func getMessage(c echo.Context) error {
	userID := sessUserID(c)
	if userID == 0 {
		return c.NoContent(http.StatusForbidden)
	}

	chanID, err := strconv.ParseInt(c.QueryParam("channel_id"), 10, 64)
	if err != nil {
		return err
	}
	lastID, err := strconv.ParseInt(c.QueryParam("last_message_id"), 10, 64)
	if err != nil {
		return err
	}

	msgWithUsrs, err := queryMsgWithUsrs(chanID, lastID)
	if err != nil {
		return err
	}

	response := make([]map[string]interface{}, 0)
	for i := len(msgWithUsrs) - 1; i >= 0; i-- {
		mwu := msgWithUsrs[i]
		r := pureJsonifyMessage(mwu)
		response = append(response, r)
	}

	// messages はid降順
	if len(msgWithUsrs) > 0 {
		_, err := db.Exec("INSERT INTO haveread (user_id, channel_id, message_id, updated_at, created_at)"+
			" VALUES (?, ?, ?, NOW(), NOW())"+
			" ON DUPLICATE KEY UPDATE message_id = ?, updated_at = NOW()",
			userID, chanID, msgWithUsrs[0].Msg.ID, msgWithUsrs[0].Msg.ID)
		if err != nil {
			return err
		}
	}

	return c.JSON(http.StatusOK, response)
}

func addMessage(channelID, userID int64, content string) (int64, error) {
	res, err := db.Exec(
		"INSERT INTO message (channel_id, user_id, content, created_at) VALUES (?, ?, ?, NOW())",
		channelID, userID, content)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
