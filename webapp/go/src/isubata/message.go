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

func queryMessagesWithJsonify(chanID, lastID int64) ([]map[string]interface{}, error) {
	response := make([]map[string]interface{}, 100)

	query := `SELECT u.name, u.display_name, u.avatar_icon, m.id, m.created_at, m.content FROM user u INNER JOIN message m ON u.id = m.user_id WHERE m.id > ? AND m.channel_id = ? ORDER BY m.id ASC LIMIT 100`
	rows, err := db.Query(query, lastID, chanID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		r := make(map[string]interface{})

		u := User{}
		m := Message{}

		rows.Scan(&u.Name, &u.DisplayName, &u.AvatarIcon, &m.ID, &m.CreatedAt, &m.Content)
		if err != nil {
			return nil, err
		}
		r["id"] = m.ID
		r["user"] = u
		r["date"] = m.CreatedAt.Format("2006/01/02 15:04:05")
		r["content"] = m.Content

		response = append(response, r)
	}

	return response, nil
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

	response, err := queryMessagesWithJsonify(chanID, lastID)
	if err != nil {
		return err
	}

	if len(response) > 0 {
		_, err := db.Exec("INSERT INTO haveread (user_id, channel_id, message_id, updated_at, created_at)"+
			" VALUES (?, ?, ?, NOW(), NOW())"+
			" ON DUPLICATE KEY UPDATE message_id = ?, updated_at = NOW()",
			userID, chanID, response[0]["id"], response[0]["id"])
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
