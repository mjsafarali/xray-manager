package v1

import (
	"encoding/base64"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"net/http"
	"shadowsocks-manager/internal/coordinator"
	"shadowsocks-manager/internal/database"
	"shadowsocks-manager/internal/utils"
)

type ProfileResponse struct {
	User            database.User `json:"user"`
	Used            float64       `json:"used"`
	ShadowsocksLink string        `json:"shadowsocks_link"`
}

func ProfileShow(d *database.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user *database.User
		for _, u := range d.Data.Users {
			if u.Identity == c.QueryParam("u") {
				user = u
			}
		}
		if user == nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Not found.",
			})
		}

		s := d.Data.Settings
		auth := base64.StdEncoding.EncodeToString([]byte(user.Method + ":" + user.Password))
		link := fmt.Sprintf("ss://%s@%s:%d#%s", auth, s.ShadowsocksHost, s.ShadowsocksPort, user.Name)
		used := utils.RoundFloat(user.Used*d.Data.Settings.TrafficRatio, 2)

		r := ProfileResponse{User: *user, ShadowsocksLink: link, Used: used}
		r.User.Used = 0
		r.User.Quota = int(float64(r.User.Quota) * d.Data.Settings.TrafficRatio)

		return c.JSON(http.StatusOK, r)
	}
}

func ProfileReset(coordinator *coordinator.Coordinator, d *database.Database) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user *database.User
		for _, u := range d.Data.Users {
			if u.Identity == c.QueryParam("u") {
				user = u
			}
		}
		if user == nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Not found.",
			})
		}

		var newPassword string
		for {
			newPassword = random.String(16)
			isUnique := true
			for _, u := range d.Data.Users {
				if u.Password == newPassword {
					isUnique = false
					break
				}
			}
			if isUnique {
				break
			}
		}
		user.Password = newPassword
		d.Save()

		go coordinator.SyncUsers()

		return c.JSON(http.StatusOK, user)
	}
}
