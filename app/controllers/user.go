package controllers

import (
	"notification-ms/app"

	"github.com/revel/revel"
)

type User struct {
	*revel.Controller
}

func getUserDeviceRedis(id string) string {
	userDev, errud := app.RedisClient.Get("userdev/" + id).Result()
	if errud != nil {
		userDev = "No device"
	}
	return userDev
}

func setUserDeviceRedis(id string, dev string) {
	app.RedisClient.Set("userdev/"+id, dev, 0).Result()
}

func (ud User) Show() revel.Result {
	data := make(map[string]interface{})
	userId := ud.Params.Route.Get("user")
	data["user"] = userId
	data["device"] = getUserDeviceRedis(userId)
	return ud.RenderJSON(data)
}

func (ud User) Create() revel.Result {
	data := bindDataSafeUd(ud)
	if data["error"] == "errParams" {
		return ud.RenderJSON(data)
	} else {
		setUserDeviceRedis(data["user"].(string), data["device"].(string))
		data["device"] = getUserDeviceRedis(data["user"].(string))
		return ud.RenderJSON(data)
	}
}

func bindDataSafeUd(ud User) map[string]interface{} {
	var data map[string]interface{}
	ud.Params.BindJSON(&data)
	if data == nil {
		data = make(map[string]interface{})
	}
	if data["user"] == nil || data["device"] == nil {
		data["error"] = "errParams"
	}
	return data
}
