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
	var data = make(map[string]interface{})
	data["user"], _ = app.RedisClient.Get("notif/" + id + "/user").Result()
	data["device"], _ = app.RedisClient.Get("notif/" + id + "/device").Result()
	data["data"], _ = app.RedisClient.Get("notif/" + id + "/data").Result()
	return data
}

func storeNotification(data map[string]interface{}) string {
	id := strconv.Itoa(rand.Intn(1000))
	app.RedisClient.Set("notif/"+id+"/user", data["user"], 0).Result()
	app.RedisClient.Set("notif/"+id+"/device", data["device"], 0).Result()
	app.RedisClient.Set("notif/"+id+"/data", data, 0).Result()
	return id
}

func (notif Notification) Show() revel.Result {
	return notif.RenderJSON(getNotificationRedisDic(notif.Params.Route.Get("id")))
}

func (notif Notification) Send() revel.Result {
	JSON := getJSON(notif)
	JSON, err := reviewJSON(JSON)
	if err {
		return notif.RenderJSON(JSON)
	}
	if sendFirebaseNotification(JSON) == false {
		return notif.RenderJSON("Error sending notification to Firebase")
	}
	id := storeNotification(JSON)
	returnData := getNotificationRedisDic(id)
	return notif.RenderJSON(returnData)
}

func getJSON(notif Notification) map[string]interface{} {
	var json map[string]interface{}
	notif.Params.BindJSON(&json)
	return json
}

func reviewJSON(data map[string]interface{}) (map[string]interface{}, bool) {
	var errors = make(map[string]interface{})
	var err = false
	if data == nil {
		data = make(map[string]interface{})
		data["error"] = errors
		errors["Parse JSON"] = "There is not a valid JSON in your Request"
		err = true
		return data, err
	}
	if data["device"] == nil && data["user"] != nil {
		newDevice := getUserDeviceRedis(data["user"].(string))
		if newDevice != "No device" {
			data["device"] = newDevice
		} else {
			errors["Device"] = "Couldn't recover device from user"
		}
	} else {
		errors["User Device"] = "We need a device to send Notifications!"
	}
	return data, err
}
