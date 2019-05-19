package controllers

import (
	"math/rand"
	"notification-ms/app"
	"strconv"

	"github.com/revel/revel"
)

type Notification struct {
	*revel.Controller
}

func getNotificationRedisDic(id string) map[string]interface{} {
	data := make(map[string]interface{})
	user, erru := app.RedisClient.Get("notif/" + id + "/user").Result()
	device, errd := app.RedisClient.Get("notif/" + id + "/device").Result()
	topic, errt := app.RedisClient.Get("notif/" + id + "/topic").Result()
	message, errm := app.RedisClient.Get("notif/" + id + "/message").Result()
	if erru != nil || errd != nil || errt != nil || errm != nil {
		data["id"] = "errGetingRedis"
		data["user"] = "errGetingRedis"
		data["device"] = "errGetingRedis"
		data["topic"] = "errGetingRedis"
		data["message"] = "errGetingRedis"
	} else {
		data["id"] = id
		data["user"] = user
		data["device"] = device
		data["topic"] = topic
		data["message"] = message
	}
	return data
}

func storeNorification(user string, device string, topic string, message string) string {
	id := strconv.Itoa(rand.Intn(1000))
	app.RedisClient.Set("notif/"+id+"/user", user, 0).Result()
	app.RedisClient.Set("notif/"+id+"/device", device, 0).Result()
	app.RedisClient.Set("notif/"+id+"/topic", topic, 0).Result()
	app.RedisClient.Set("notif/"+id+"/message", message, 0).Result()
	return id
}

func sentNotification(user string, device string, topic string, message string) {
	sendAndroidNotification(device, topic, message)
}

func (notif Notification) Show() revel.Result {
	return notif.RenderJSON(getNotificationRedisDic(notif.Params.Route.Get("id")))
}

func (notif Notification) Send() revel.Result {
	data := bindDataSafeNotif(notif)
	if data["error"] == "errParams" {
		return notif.RenderJSON(data)
	} else {
		user := data["user"].(string)
		device := data["device"].(string)
		topic := data["topic"].(string)
		message := data["message"].(string)
		id := storeNorification(user, device, topic, message)
		sentNotification(user, device, topic, message)
		data = getNotificationRedisDic(id)
		return notif.RenderJSON(data)
	}
}

func bindDataSafeNotif(notif Notification) map[string]interface{} {
	var data map[string]interface{}
	notif.Params.BindJSON(&data)
	if data == nil {
		data = make(map[string]interface{})
	}
	if data["device"] == nil && data["user"] != nil {
		newDevice := getUserDeviceRedis(data["user"].(string))
		if newDevice != "No device" {
			data["device"] = newDevice
		} else {
			data["device"] = "Couldn't recover device from user"
			data["error"] = "errParams"
		}
	}
	if data["user"] == nil || data["device"] == nil || data["topic"] == nil || data["message"] == nil {
		data["error"] = "errParams"
	}
	return data
}
