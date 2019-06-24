package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"notification-ms/app"
)

func sendFirebaseNotification(data map[string]interface{}) bool {
	url := "https://fcm.googleapis.com/fcm/send"
	dataToSend := make(map[string]interface{})
	dataToSend["data"] = data
	dataToSend["to"] = data["device"]
	dataToSend["notification"] = make(map[string]interface{})
	dataToSend["notification"].(map[string]interface{})["title"] = data["title"]
	dataToSend["notification"].(map[string]interface{})["body"] = data["body"]
	jsonStr, err := json.Marshal(dataToSend)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", app.CliFirebaseApiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var respJSON interface{}
	json.Unmarshal(body, &respJSON)
	return respJSON.(map[string]interface{})["failure"] == 0.0
}
