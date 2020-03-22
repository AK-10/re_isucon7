package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func postProfile(c echo.Context) error {
	self, err := ensureLogin(c)
	if self == nil {
		return err
	}

	avatarName := ""
	var avatarData []byte

	if fh, err := c.FormFile("avatar_icon"); err == http.ErrMissingFile {
		// no file upload
	} else if err != nil {
		return err
	} else {
		dotPos := strings.LastIndexByte(fh.Filename, '.')
		if dotPos < 0 {
			return ErrBadReqeust
		}
		ext := fh.Filename[dotPos:]
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif":
			break
		default:
			return ErrBadReqeust
		}

		file, err := fh.Open()
		if err != nil {
			return err
		}
		avatarData, _ = ioutil.ReadAll(file)
		file.Close()

		if len(avatarData) > avatarMaxBytes {
			return ErrBadReqeust
		}

		avatarName = fmt.Sprintf("%x%s", sha1.Sum(avatarData), ext)
	}

	if avatarName != "" && len(avatarData) > 0 {
		if err := writeIconFile(avatarName, avatarData); err != nil {
			return err
		}

		// _, err := db.Exec("INSERT INTO image (name, data) VALUES (?, ?)", avatarName, avatarData)
		// if err != nil {
		// 	return err
		// }
		_, err = db.Exec("UPDATE user SET avatar_icon = ? WHERE id = ?", avatarName, self.ID)
		if err != nil {
			return err
		}
	}

	if name := c.FormValue("display_name"); name != "" {
		_, err := db.Exec("UPDATE user SET display_name = ? WHERE id = ?", name, self.ID)
		if err != nil {
			return err
		}
	}

	return c.Redirect(http.StatusSeeOther, "/")
}

func writeIconFile(name string, data []byte) error {
	destPath := "/home/isucon/isubata/webapp/public/icons/"

	return ioutil.WriteFile(destPath+name, data, 0644)
	// globalCache.Set("icon-"+name, data, cache.DefaultExpiration)

	// return nil
}

func writeInitialIconFile() error {
	rows, err := db.Query("SELECT name, data FROM image")
	if err != nil {
		return err
	}

	for rows.Next() {
		var fileName string
		var data []byte

		err := rows.Scan(&fileName, &data)
		if err != nil {
			return err
		}

		// writeIconFile(fileName, data)
		if err := writeIconFile(fileName, data); err != nil {
			return err
		}
	}

	return nil
}

// fileに書き込み, nginxで配信するためいらん
// なぜかうまく動かないので，一旦go-cacheに格納
// func getIcon(c echo.Context) error {
// var name string
// var data []byte
// err := db.QueryRow("SELECT name, data FROM image WHERE name = ?",
// 	c.Param("file_name")).Scan(&name, &data)
// if err == sql.ErrNoRows {
// 	return echo.ErrNotFound
// }
// if err != nil {
// 	return err
// }
// 	name := c.Param("file_name")
// 	dInterface, ok := globalCache.Get("icon-" + name)
// 	if !ok {
// 		return errors.New("icon not found on go cache")
// 	}

// 	data, ok := dInterface.([]byte)
// 	if !ok {
// 		return errors.New("failed type assertion")
// 	}

// 	mime := ""
// 	switch true {
// 	case strings.HasSuffix(name, ".jpg"), strings.HasSuffix(name, ".jpeg"):
// 		mime = "image/jpeg"
// 	case strings.HasSuffix(name, ".png"):
// 		mime = "image/png"
// 	case strings.HasSuffix(name, ".gif"):
// 		mime = "image/gif"
// 	default:
// 		return echo.ErrNotFound
// 	}
// 	return c.Blob(http.StatusOK, mime, data)
// }
