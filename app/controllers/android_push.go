package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"notification-ms/app"

	"github.com/prometheus/common/log"
)

type answer struct {
	notifications []notification
}
type notification struct {
	tokens   []string
	platform string
	message  string
	title    string
}

func sendAndroidNotification(token string, message string, title string) {
	var answer = make(map[string]interface{})
	answer["notifications"] = []map[string]interface{}{
		make(map[string]interface{}),
	}
	answer["notifications"].([]map[string]interface{})[0]["tokens"] = []string{
		token,
	}
	answer["notifications"].([]map[string]interface{})[0]["platform"] = 2
	answer["notifications"].([]map[string]interface{})[0]["message"] = "message"
	answer["notifications"].([]map[string]interface{})[0]["title"] = "title"
	jsonValues, _ := json.Marshal(answer)
	_, err := http.Post("http://"+app.GoRushHost+"/api/push", "application/json", bytes.NewBuffer(jsonValues))
	bytes.NewBuffer(jsonValues)
	if err != nil {
		log.Error("Could not send notification")
	}
}
